@echo off
rem Выгружает все модели из видеопамяти. Сама Ollama продолжает работать,
rem следующий запрос к модели снова поднимет её в память (~30-60 сек).
for /f "skip=1 tokens=1" %%m in ('ollama ps') do ollama stop %%m
echo.
ollama ps
pause
