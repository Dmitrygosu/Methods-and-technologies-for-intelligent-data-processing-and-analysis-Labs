# Небольшой набор функций для работы с текстом — цель для генератора тестов.

def word_count(text):
    """Возвращает количество слов в строке.

    Слова разделяются любым количеством пробельных символов.
    Для пустой строки возвращает 0.
    Если text не строка — бросает TypeError.

    >>> word_count("привет мир")
    2
    >>> word_count("")
    0
    """
    if not isinstance(text, str):
        raise TypeError("ожидалась строка")
    return len(text.split())


def is_palindrome(text):
    """Проверяет, является ли строка палиндромом.

    Регистр и пробелы не учитываются.
    Пустая строка считается палиндромом.

    >>> is_palindrome("А роза упала на лапу Азора")
    True
    >>> is_palindrome("привет")
    False
    """
    cleaned = "".join(text.lower().split())
    return cleaned == cleaned[::-1]


def truncate(text, limit):
    """Обрезает строку до limit символов.

    Если строка длиннее limit, обрезает и добавляет многоточие "...",
    при этом итоговая длина равна ровно limit.
    Если строка помещается — возвращает её без изменений.
    Если limit < 4 — бросает ValueError, т.к. места не хватит даже на многоточие.

    >>> truncate("привет мир", 8)
    'приве...'
    >>> truncate("да", 10)
    'да'
    """
    if limit < 4:
        raise ValueError("limit должен быть не меньше 4")
    if len(text) <= limit:
        return text
    return text[:limit - 3] + "..."


def normalize_spaces(text):
    """Убирает лишние пробелы: по краям и повторяющиеся внутри.

    Если text не строка — бросает TypeError.

    >>> normalize_spaces("  привет   мир  ")
    'привет мир'
    """
    if not isinstance(text, str):
        raise TypeError("ожидалась строка")
    return " ".join(text.split())


def safe_div(a, b, default=0.0):
    """Делит a на b. При делении на ноль возвращает default.

    >>> safe_div(10, 4)
    2.5
    >>> safe_div(1, 0)
    0.0
    >>> safe_div(1, 0, default=-1)
    -1
    """
    if b == 0:
        return default
    return a / b
