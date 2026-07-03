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
