@echo off
ollama pull qwen2.5:3b
ollama create pochta-analitik -f "%~dp0pochta-analitik.Modelfile"
ollama list
