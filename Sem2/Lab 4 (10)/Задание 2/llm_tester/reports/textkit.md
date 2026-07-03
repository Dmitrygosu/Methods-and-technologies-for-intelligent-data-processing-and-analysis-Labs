# Отчёт по тестированию textkit.py

Модель-генератор: test-inzhener, попыток генерации: 1, кругов доработки: 1

## Динамика по кругам

| круг | всего | прошло | упало | ошибок |
|------|-------|--------|-------|--------|
| 1 | 9 | 9 | 0 | 0 |

## Статистика финального прогона

- всего тестов: 9
- прошло: 9
- упало: 0
- ошибок: 0
- время: 0.03 с

## Анализ модели

Прогон автотестов для модуля `textkit.py` завершен успешно. Всего было выполнено 9 тестов, из которых все прошли без ошибок. Однако стоит отметить, что не были проверены крайние случаи для функций `truncate` и `normalize_spaces`. Например, `truncate` не был прощен с ситуацией, где длина исходной строки меньше максимального допустимого значения. Также не было тестов на случай, если в функцию `normalize_spaces` передается строка без пробелов или только пробельных символов.

## Сгенерированные тесты

```python
import pytest

import textkit

def test_word_count():
    assert textkit.word_count("привет мир") == 2
    assert textkit.word_count("") == 0
    
def test_word_count_type_error():
    with pytest.raises(TypeError):
        textkit.word_count(123)
    
def test_is_palindrome():
    assert textkit.is_palindrome("А роза упала на лапу Азора") == True
    assert textkit.is_palindrome("привет") == False
    
def test_is_palindrome_empty_string():
    assert textkit.is_palindrome("") == True

def test_truncate_normal_case():
    assert textkit.truncate("привет мир", 8) == "приве..."
    
def test_truncate_with_limit_less_than_4():
    with pytest.raises(ValueError):
        textkit.truncate("привет мир", 3)
    
def test_normalize_spaces_normal_case():
    assert textkit.normalize_spaces("  привет   мир  ") == "привет мир"
    
def test_safe_div_normal_case():
    assert textkit.safe_div(10, 4) == 2.5
    assert textkit.safe_div(1, 0) == 0.0
    
def test_safe_div_default_argument():
    assert textkit.safe_div(1, 0, default=-1) == -1
```

## Лог pytest

```
============================= test session starts =============================
platform win32 -- Python 3.14.3, pytest-9.1.1, pluggy-1.6.0 -- C:\Python314\python.exe
cachedir: .pytest_cache
rootdir: C:\Users\ignat\Desktop\shish\Methods-and-technologies-for-intelligent-data-processing-and-analysis-Labs\Sem2\Lab 4 (10)\llm_tester
collecting ... collected 9 items

generated/test_textkit.py::test_word_count PASSED                        [ 11%]
generated/test_textkit.py::test_word_count_type_error PASSED             [ 22%]
generated/test_textkit.py::test_is_palindrome PASSED                     [ 33%]
generated/test_textkit.py::test_is_palindrome_empty_string PASSED        [ 44%]
generated/test_textkit.py::test_truncate_normal_case PASSED              [ 55%]
generated/test_textkit.py::test_truncate_with_limit_less_than_4 PASSED   [ 66%]
generated/test_textkit.py::test_normalize_spaces_normal_case PASSED      [ 77%]
generated/test_textkit.py::test_safe_div_normal_case PASSED              [ 88%]
generated/test_textkit.py::test_safe_div_default_argument PASSED         [100%]

- generated xml file: C:\Users\ignat\Desktop\shish\Methods-and-technologies-for-intelligent-data-processing-and-analysis-Labs\Sem2\Lab 4 (10)\llm_tester\generated\test_textkit.py.junit.xml -
============================== 9 passed in 0.08s ==============================
```
