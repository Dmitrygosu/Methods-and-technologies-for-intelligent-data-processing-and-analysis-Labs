SHELL := sh
LAB := $(CURDIR)
# ollama ставится в профиль пользователя; если её нет в PATH текущей
# оболочки — берём по абсолютному пути
OLLAMA := $(if $(shell command -v ollama 2>/dev/null),ollama,"$(LOCALAPPDATA)/Programs/Ollama/ollama.exe")

.PHONY: all models task1 task2 task3 task4 check n8n-up n8n-stop \
        llm-stop llm-kill llm-start status start stop

all: models

# собрать все шесть специализированных моделей
models:
	$(OLLAMA) pull qwen2.5:3b
	cd "$(LAB)/Задание 1/models" && $(OLLAMA) create meteo-sovetnik -f meteo-sovetnik.Modelfile
	cd "$(LAB)/Задание 2/models" && $(OLLAMA) create test-inzhener -f test-inzhener.Modelfile
	cd "$(LAB)/Задание 3/models" && $(OLLAMA) create optimist -f optimist.Modelfile && \
		$(OLLAMA) create skeptik -f skeptik.Modelfile && $(OLLAMA) create ekspert -f ekspert.Modelfile
	cd "$(LAB)/Задание 4/models" && $(OLLAMA) create pochta-analitik -f pochta-analitik.Modelfile
	$(OLLAMA) list

# ----- запуск заданий -----

# погодный советник: http://127.0.0.1:5000
task1:
	cd "$(LAB)/Задание 1/weather_advisor" && python app.py

# генерация и прогон тестов по обеим целям
task2:
	cd "$(LAB)/Задание 2/llm_tester" && python testgen.py targets/textkit.py
	cd "$(LAB)/Задание 2/llm_tester" && python testgen.py targets/statkit.py

# чат с тремя ролями и почтовый ассистент живут в n8n
task3 task4: n8n-up
	@echo "Открой http://localhost:5678 и запусти нужный воркфлоу"

# автопроверка советника: сверяет решение про зонт с правилом
# на контрольных городах, печатает таблицу вида таблицы 2 из отчёта
check:
	cd "$(LAB)/Задание 1/weather_advisor" && python check_advisor.py

# ----- запуск/остановка всего стенда одной командой -----

start:
	powershell -NoProfile -ExecutionPolicy Bypass -File "$(LAB)/start_all.ps1"

stop:
	powershell -NoProfile -ExecutionPolicy Bypass -File "$(LAB)/stop_all.ps1"

# ----- n8n -----

n8n-up:
	-docker start n8n 2>/dev/null || docker run -d --name n8n -p 5678:5678 \
		-e N8N_SECURE_COOKIE=false -e GENERIC_TIMEZONE=Europe/Ulyanovsk \
		-e N8N_RESTRICT_FILE_ACCESS_TO=/data/output \
		-v n8n_data:/home/node/.n8n \
		-v "$(LAB)/Задание 4/n8n/output:/data/output" n8nio/n8n:latest
	docker ps --filter name=n8n --format "{{.Names}} {{.Status}}"

n8n-stop:
	docker stop n8n

# ----- управление ллмками -----

# выгрузить все модели из видеопамяти (Ollama остаётся работать)
llm-stop:
	-$(OLLAMA) ps | tail -n +2 | awk '{print $$1}' | xargs -r -n1 $(OLLAMA) stop
	$(OLLAMA) ps

# полностью вырубить Ollama
llm-kill:
	-taskkill //IM "ollama app.exe" //F
	-taskkill //IM ollama.exe //F

# поднять Ollama обратно
llm-start:
	cmd //c start "" "$(LOCALAPPDATA)\Programs\Ollama\ollama app.exe"

# что сейчас запущено и сколько памяти ест
status:
	-$(OLLAMA) ps
	-nvidia-smi --query-gpu=memory.used,memory.total --format=csv,noheader || true
	-docker ps --filter name=n8n --format "n8n: {{.Status}}"
