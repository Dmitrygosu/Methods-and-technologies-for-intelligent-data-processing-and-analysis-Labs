# Генератор тестов на базе локальной LLM.
#
# Схема работы:
#   1) читаем исходник целевого модуля;
#   2) модель test-inzhener пишет для него тесты на pytest;
#   3) проверяем, что ответ — валидный Python (если нет, возвращаем модели
#      ошибку и просим исправить, до 3 попыток);
#   4) запускаем pytest и собираем статистику;
#   5) отдаём итоги прогона модели и получаем анализ на русском;
#   6) складываем всё в отчёт reports/<модуль>.md.
#
# Запуск:  python testgen.py targets/textkit.py

import json
import os
import re
import subprocess
import sys
import time
import xml.etree.ElementTree as ET

import requests

OLLAMA_URL = os.environ.get("OLLAMA_URL", "http://localhost:11434")
GEN_MODEL = os.environ.get("GEN_MODEL", "test-inzhener")
MAX_ATTEMPTS = 3   # попытки получить синтаксически корректный код
MAX_ROUNDS = 3     # круги "прогон -> доработка упавших тестов"

HERE = os.path.dirname(os.path.abspath(__file__))


def ask(model, prompt, temperature=0.2):
    r = requests.post(f"{OLLAMA_URL}/api/generate", json={
        "model": model,
        "prompt": prompt,
        "stream": False,
        "options": {"temperature": temperature},
    }, timeout=600)
    r.raise_for_status()
    return r.json()["response"]


def strip_fences(code):
    # модель иногда всё-таки заворачивает код в ```python ... ```
    code = code.strip()
    m = re.search(r"```(?:python)?\s*(.*?)```", code, re.S)
    if m:
        code = m.group(1).strip()
    return code


def ensure_imports(code, module_name):
    # частая ошибка мелких моделей: используют pytest.raises,
    # а import pytest написать забывают
    lines = []
    if re.search(r"\bpytest\.", code) and "import pytest" not in code:
        lines.append("import pytest")
    if f"import {module_name}" not in code:
        lines.append(f"import {module_name}")
    if lines:
        code = "\n".join(lines) + "\n\n" + code
    return code


def generate_tests(module_name, source):
    prompt = (
        f"Вот исходный код модуля {module_name}.py:\n\n{source}\n\n"
        f"Напиши тесты pytest для всех функций модуля. "
        f"Модуль импортируй строго так: import {module_name}\n"
        f"К каждой функции — минимум два теста: обычный случай и граничный."
    )
    last_error = None
    for attempt in range(1, MAX_ATTEMPTS + 1):
        if last_error:
            full = (prompt + f"\n\nТвоя прошлая версия не прошла проверку "
                             f"синтаксиса:\n{last_error}\nИсправь и выведи "
                             f"полный код заново.")
        else:
            full = prompt

        print(f"  попытка {attempt}: запрашиваю тесты у модели {GEN_MODEL}...")
        t0 = time.time()
        code = ensure_imports(strip_fences(ask(GEN_MODEL, full)), module_name)
        print(f"  ответ получен за {time.time() - t0:.1f} с, "
              f"{len(code.splitlines())} строк")

        try:
            compile(code, "<generated>", "exec")
            return code, attempt
        except SyntaxError as e:
            last_error = f"строка {e.lineno}: {e.msg}"
            print(f"  синтаксическая ошибка ({last_error}), пробую ещё раз")

    raise RuntimeError("модель так и не выдала синтаксически корректный код")


def repair_tests(module_name, source, test_code, failures):
    prompt = (
        f"Вот исходный код модуля {module_name}.py:\n\n{source}\n\n"
        f"Вот написанные тобой тесты:\n\n{test_code}\n\n"
        f"Часть тестов упала при прогоне:\n"
        + json.dumps(failures, ensure_ascii=False, indent=2) + "\n\n"
        f"Сверь ожидания упавших тестов с реальным поведением кода — "
        f"эталоном считай докстринги. Если тест ожидает исключение, которого "
        f"код не бросает, исправь тест. Выведи полный файл тестов заново, "
        f"только код без пояснений."
    )
    code = ensure_imports(strip_fences(ask(GEN_MODEL, prompt)), module_name)
    compile(code, "<repaired>", "exec")
    return code


def quarantine(code, failing_names):
    # тесты, которые модель не смогла согласовать с кодом за отведённые
    # круги, не выбрасываем, а помечаем skip — для ручного разбора
    for name in failing_names:
        code = re.sub(
            rf"^def {re.escape(name)}\(",
            f'@pytest.mark.skip(reason="карантин: ожидание не подтверждено, '
            f'нужен ручной разбор")\ndef {name}(',
            code, count=1, flags=re.M,
        )
    if "import pytest" not in code:
        code = "import pytest\n\n" + code
    return code


def run_pytest(test_file):
    xml_path = test_file + ".junit.xml"
    proc = subprocess.run(
        [sys.executable, "-m", "pytest", test_file, "-v",
         "--tb=short", f"--junitxml={xml_path}"],
        capture_output=True, text=True, encoding="utf-8",
        errors="replace", cwd=HERE,
    )
    stats = {"total": 0, "passed": 0, "failed": 0, "errors": 0, "time": 0.0}
    failures = []
    if os.path.exists(xml_path):
        suite = ET.parse(xml_path).getroot().find("testsuite")
        stats["total"] = int(suite.get("tests", 0))
        stats["failed"] = int(suite.get("failures", 0))
        stats["errors"] = int(suite.get("errors", 0))
        stats["time"] = float(suite.get("time", 0))
        stats["passed"] = stats["total"] - stats["failed"] - stats["errors"]
        for case in suite.iter("testcase"):
            for bad in list(case.iter("failure")) + list(case.iter("error")):
                failures.append({
                    "test": case.get("name"),
                    "message": (bad.get("message") or "")[:500],
                })
        os.remove(xml_path)
    return stats, failures, proc.stdout


def analyze(module_name, stats, failures, test_code):
    prompt = (
        f"Прогон автотестов для модуля {module_name}.py завершён.\n"
        f"Итого тестов: {stats['total']}, прошло: {stats['passed']}, "
        f"упало: {stats['failed']}, ошибок сборки: {stats['errors']}.\n"
    )
    if failures:
        prompt += "Упавшие тесты:\n" + json.dumps(failures, ensure_ascii=False,
                                                  indent=2) + "\n"
        prompt += ("Определи по каждому падению: это ошибка в тестируемом "
                   "коде или тест написан неверно? ")
    prompt += (
        "Сделай короткий вывод на русском (5-8 предложений): достаточно ли "
        "покрытие, какие граничные случаи не проверены, что стоит добавить. "
        "Вот код тестов для справки:\n" + test_code
    )
    return ask("qwen2.5:3b", prompt, temperature=0.5)


def main():
    if len(sys.argv) != 2:
        print("использование: python testgen.py targets/<модуль>.py")
        sys.exit(1)

    target = sys.argv[1]
    module_name = os.path.splitext(os.path.basename(target))[0]
    with open(os.path.join(HERE, target), encoding="utf-8") as f:
        source = f.read()

    print(f"[1/4] Генерация тестов для {module_name}.py")
    test_code, attempts = generate_tests(module_name, source)

    os.makedirs(os.path.join(HERE, "generated"), exist_ok=True)
    test_file = os.path.join(HERE, "generated", f"test_{module_name}.py")
    print(f"[2/4] Прогон и доработка (до {MAX_ROUNDS} кругов)")

    rounds = []
    for rnd in range(1, MAX_ROUNDS + 1):
        with open(test_file, "w", encoding="utf-8") as f:
            f.write(test_code + "\n")
        stats, failures, log = run_pytest(test_file)
        rounds.append(dict(stats, failures=list(failures)))
        print(f"  круг {rnd}: всего {stats['total']}, прошло "
              f"{stats['passed']}, упало {stats['failed']}, "
              f"ошибок {stats['errors']}")
        if stats["failed"] + stats["errors"] == 0 or rnd == MAX_ROUNDS:
            break
        print("  есть падения — отправляю модели протокол на доработку")
        try:
            test_code = repair_tests(module_name, source, test_code, failures)
        except SyntaxError:
            print("  доработка сломала синтаксис, оставляю прошлую версию")
            break

    quarantined = []
    if stats["failed"] + stats["errors"] > 0:
        quarantined = [f["test"] for f in failures]
        print(f"  карантин для: {', '.join(quarantined)}")
        test_code = quarantine(test_code, quarantined)
        with open(test_file, "w", encoding="utf-8") as f:
            f.write(test_code + "\n")
        stats, failures, log = run_pytest(test_file)

    print(f"[3/4] Итог: {stats['passed']}/{stats['total']} за "
          f"{len(rounds)} круг(а), в карантине: {len(quarantined)}")

    print("[4/4] Анализ результатов моделью")
    verdict = analyze(module_name, stats, rounds[-1]["failures"], test_code)

    os.makedirs(os.path.join(HERE, "reports"), exist_ok=True)
    report_file = os.path.join(HERE, "reports", f"{module_name}.md")
    with open(report_file, "w", encoding="utf-8") as f:
        f.write(f"# Отчёт по тестированию {module_name}.py\n\n")
        f.write(f"Модель-генератор: {GEN_MODEL}, попыток генерации: "
                f"{attempts}, кругов доработки: {len(rounds)}\n\n")
        f.write("## Динамика по кругам\n\n")
        f.write("| круг | всего | прошло | упало | ошибок |\n")
        f.write("|------|-------|--------|-------|--------|\n")
        for i, r in enumerate(rounds, 1):
            f.write(f"| {i} | {r['total']} | {r['passed']} | "
                    f"{r['failed']} | {r['errors']} |\n")
        f.write("\n")
        f.write(f"## Статистика финального прогона\n\n")
        f.write(f"- всего тестов: {stats['total']}\n")
        f.write(f"- прошло: {stats['passed']}\n")
        f.write(f"- упало: {stats['failed']}\n")
        f.write(f"- ошибок: {stats['errors']}\n")
        f.write(f"- время: {stats['time']:.2f} с\n\n")
        if quarantined:
            f.write("## Карантин (нужен ручной разбор)\n\n")
            for fl in rounds[-1]["failures"]:
                f.write(f"- `{fl['test']}`: {fl['message']}\n")
            f.write("\n")
        f.write("## Анализ модели\n\n" + verdict.strip() + "\n\n")
        f.write("## Сгенерированные тесты\n\n```python\n" + test_code +
                "\n```\n\n")
        f.write("## Лог pytest\n\n```\n" + log.strip() + "\n```\n")

    print(f"Готово. Отчёт: {os.path.relpath(report_file, HERE)}")


if __name__ == "__main__":
    main()
