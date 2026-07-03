@echo off
rem Поднимает Ollama обратно после llm_kill.cmd
start "" "%LOCALAPPDATA%\Programs\Ollama\ollama app.exe"
timeout /t 4 /nobreak >nul
ollama list
pause
