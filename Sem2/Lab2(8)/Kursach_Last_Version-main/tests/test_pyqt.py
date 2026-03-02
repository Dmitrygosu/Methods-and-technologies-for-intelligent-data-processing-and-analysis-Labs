from PyQt6.QtWidgets import QApplication, QLabel
import sys
try:
    app = QApplication(sys.argv)
    label = QLabel("Test")
    print("PyQt6 imported and initialized successfully")
except Exception as e:
    print(f"Failed to initialize PyQt6: {e}")
