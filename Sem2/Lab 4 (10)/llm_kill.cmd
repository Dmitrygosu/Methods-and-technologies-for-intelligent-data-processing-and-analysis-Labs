@echo off
rem Полностью вырубает Ollama (и модели, и фоновый процесс).
rem Обратно: llm_on.cmd или просто запустить Ollama из меню Пуск.
taskkill /IM "ollama app.exe" /F 2>nul
taskkill /IM ollama.exe /F 2>nul
echo Ollama остановлена.
pause
