@echo off
ollama pull qwen2.5:3b
ollama create test-inzhener -f "%~dp0test-inzhener.Modelfile"
ollama list
