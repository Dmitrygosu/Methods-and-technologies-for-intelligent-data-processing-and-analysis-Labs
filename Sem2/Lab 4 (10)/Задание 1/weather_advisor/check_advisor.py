# Автоматическая проверка советника: гоняет набор городов с разной погодой,
# сверяет решение модели про зонт с детерминированным правилом и проверяет
# список одежды на отсутствие мусора. Печатает табличку (как таблица 2
# в отчёте) и завершается ненулевым кодом, если найдено расхождение.
#
# Запуск: python check_advisor.py [город1 город2 ...]
# Без аргументов берёт контрольный набор из отчёта.

import sys
import time

import weather
import advisor

DEFAULT_CITIES = [
    "Ульяновск", "Санкт-Петербург", "Мурманск", "Сочи", "Казань",
]


def check_city(name):
    city = weather.find_city(name)
    wx = weather.get_forecast(city["lat"], city["lon"])
    t0 = time.time()
    adv = advisor.get_advice(f"{city['name']} ({city['country']})", wx)
    dt = time.time() - t0

    need = advisor.umbrella_needed(wx)
    umbrella_match = adv["umbrella"] == need
    clothing_match = advisor.clothing_ok(adv["clothing"])

    return {
        "city": city["name"],
        "temp": wx["now"]["temp"],
        "rain_prob": wx["today"]["rain_prob"],
        "model_umbrella": adv["umbrella"],
        "rule_umbrella": need,
        "umbrella_match": umbrella_match,
        "clothing_match": clothing_match,
        "clothing": adv["clothing"],
        "time": dt,
    }


def main():
    cities = sys.argv[1:] or DEFAULT_CITIES
    rows = []
    failed = False

    print(f"Проверяю {len(cities)} город(ов)...\n")
    for name in cities:
        print(f"  {name}...", end=" ", flush=True)
        try:
            row = check_city(name)
        except Exception as e:
            print(f"ОШИБКА: {e}")
            failed = True
            continue
        ok = row["umbrella_match"] and row["clothing_match"]
        print("OK" if ok else "РАСХОЖДЕНИЕ", f"({row['time']:.1f} с)")
        rows.append(row)
        if not ok:
            failed = True

    print("\nГород            | Темп | Осадки% | Модель      | Правило     | Одежда")
    print("-" * 78)
    for r in rows:
        model_txt = "брать" if r["model_umbrella"] else "не нужен"
        rule_txt = "брать" if r["rule_umbrella"] else "не нужен"
        u_flag = "  " if r["umbrella_match"] else "!!"
        c_flag = "  " if r["clothing_match"] else "!!"
        print(f"{r['city']:<17}| {r['temp']:>4} | {r['rain_prob']:>6}% | "
              f"{model_txt:<11}{u_flag}| {rule_txt:<11}{u_flag}| "
              f"{'ок' if r['clothing_match'] else 'мусор'}{c_flag}")

    bad_clothing = [r for r in rows if not r["clothing_match"]]
    if bad_clothing:
        print("\nСписки одежды с расхождением:")
        for r in bad_clothing:
            print(f"  {r['city']}: {r['clothing']}")

    print()
    if failed:
        print("ИТОГ: есть расхождения, см. выше.")
        sys.exit(1)
    else:
        print("ИТОГ: все проверки пройдены.")
        sys.exit(0)


if __name__ == "__main__":
    main()
