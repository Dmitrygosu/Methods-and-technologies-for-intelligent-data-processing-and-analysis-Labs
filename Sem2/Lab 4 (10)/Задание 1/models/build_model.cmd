@echo off
rem Сборка модели-синоптика поверх базовой qwen2.5:3b
ollama pull qwen2.5:3b
ollama create meteo-sovetnik -f "%~dp0meteo-sovetnik.Modelfile"
ollama list
