@echo off
rem Погодный советник: http://127.0.0.1:5000
cd /d "%~dp0weather_advisor"
start "" http://127.0.0.1:5000
python app.py
