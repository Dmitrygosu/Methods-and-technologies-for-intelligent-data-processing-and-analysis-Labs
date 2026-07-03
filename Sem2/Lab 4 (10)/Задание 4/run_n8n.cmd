@echo off
rem Поднимает n8n с папкой вывода дайджестов. Если контейнер уже создан
rem без этой папки (например скриптом из задания 3) — удалить его командой
rem docker rm -f n8n и запустить этот скрипт заново (воркфлоу хранятся
rem в томе n8n_data и не потеряются).
docker start n8n 2>nul || docker run -d --name n8n -p 5678:5678 ^
  -e N8N_SECURE_COOKIE=false -e GENERIC_TIMEZONE=Europe/Ulyanovsk ^
  -e N8N_RESTRICT_FILE_ACCESS_TO=/data/output ^
  -v n8n_data:/home/node/.n8n ^
  -v "%~dp0n8n\output:/data/output" n8nio/n8n:latest
start "" http://localhost:5678
