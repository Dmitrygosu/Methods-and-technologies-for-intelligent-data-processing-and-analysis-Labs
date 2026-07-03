@echo off
rem Генерация и прогон тестов для обеих целей
cd /d "%~dp0llm_tester"
python testgen.py targets/textkit.py
echo.
python testgen.py targets/statkit.py
pause
