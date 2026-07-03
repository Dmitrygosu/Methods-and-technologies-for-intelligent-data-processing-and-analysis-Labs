@echo off
rem Поднимает n8n в докере (нужен запущенный Docker Desktop)
docker start n8n 2>nul || docker run -d --name n8n -p 5678:5678 ^
  -e N8N_SECURE_COOKIE=false -e GENERIC_TIMEZONE=Europe/Ulyanovsk ^
  -v n8n_data:/home/node/.n8n n8nio/n8n:latest
start "" http://localhost:5678
