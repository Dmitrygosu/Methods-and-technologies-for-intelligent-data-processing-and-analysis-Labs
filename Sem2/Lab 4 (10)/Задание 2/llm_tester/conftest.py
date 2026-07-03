import os
import sys

# чтобы сгенерированные тесты могли писать "import textkit" без пакетов
sys.path.insert(0, os.path.join(os.path.dirname(__file__), "targets"))
