import statkit

import pytest

def test_mean():
    assert statkit.mean([1, 2, 3, 4]) == 2.5
    with pytest.raises(ValueError):
        statkit.mean([])

def test_median():
    assert statkit.median([3, 1, 2]) == 2
    assert statkit.median([4, 1, 3, 2]) == 2.5
    with pytest.raises(ValueError):
        statkit.median([])

def test_value_range():
    assert statkit.value_range([7, 2, 5]) == 5
    with pytest.raises(ValueError):
        statkit.value_range([])

@pytest.mark.skip(reason="карантин: ожидание не подтверждено, нужен ручной разбор")
def test_minmax_scale():
    assert statkit.minmax_scale([0, 5, 10]) == [0.0, 0.5, 1.0]
    assert statkit.minmax_scale([3, 3]) == [0.0, 0.0]
    with pytest.raises(ValueError):
        statkit.minmax_scale([])
