# Простейшая описательная статистика — вторая цель для генератора тестов.

def mean(values):
    """Среднее арифметическое списка чисел.

    Для пустого списка бросает ValueError.

    >>> mean([1, 2, 3, 4])
    2.5
    """
    if not values:
        raise ValueError("пустой список")
    return sum(values) / len(values)


def median(values):
    """Медиана списка чисел.

    Для чётного количества элементов — среднее двух центральных.
    Для пустого списка бросает ValueError.

    >>> median([3, 1, 2])
    2
    >>> median([4, 1, 3, 2])
    2.5
    """
    if not values:
        raise ValueError("пустой список")
    s = sorted(values)
    n = len(s)
    mid = n // 2
    if n % 2 == 1:
        return s[mid]
    return (s[mid - 1] + s[mid]) / 2


def value_range(values):
    """Размах выборки: разница между максимумом и минимумом.

    Для пустого списка бросает ValueError.

    >>> value_range([7, 2, 5])
    5
    """
    if not values:
        raise ValueError("пустой список")
    return max(values) - min(values)


def minmax_scale(values):
    """Нормирует список чисел в диапазон [0, 1].

    Если все значения одинаковые, возвращает список нулей.
    Для пустого списка возвращает пустой список.

    >>> minmax_scale([0, 5, 10])
    [0.0, 0.5, 1.0]
    >>> minmax_scale([3, 3])
    [0.0, 0.0]
    """
    if not values:
        return []
    lo, hi = min(values), max(values)
    if hi == lo:
        return [0.0 for _ in values]
    return [(v - lo) / (hi - lo) for v in values]
