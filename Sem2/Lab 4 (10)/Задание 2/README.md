# Задание 2. Система тестирования на базе локальной LLM

Модель test-inzhener генерирует pytest-тесты для выбранного модуля,
запускает их, дорабатывает упавшие по протоколу прогона (до 3 кругов),
несогласуемые тесты уводит в карантин и в конце анализирует итоги.
Отчёты складываются в llm_tester/reports/*.md.

Подготовка: `models\build_model.cmd`, затем `pip install pytest requests`.

Запуск: `run_tests.cmd` или вручную:

```
cd llm_tester
python testgen.py targets/textkit.py
python testgen.py targets/statkit.py
```
