@echo off
ollama pull qwen2.5:3b
for %%m in (optimist skeptik ekspert) do ollama create %%m -f "%~dp0%%m.Modelfile"
ollama list
