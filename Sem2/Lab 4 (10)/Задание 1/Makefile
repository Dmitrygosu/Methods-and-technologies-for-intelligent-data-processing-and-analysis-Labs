PY = python

.PHONY: setup model run exe clean

setup:
	$(PY) -m pip install -r requirements.txt

model:
	cd models && ollama pull qwen2.5:3b && \
	ollama create meteo-sovetnik -f meteo-sovetnik.Modelfile

run:
	cd weather_advisor && $(PY) app.py

exe:
	cd weather_advisor && $(PY) -m PyInstaller --onefile --noconfirm \
		--add-data "templates;templates" --name weather_advisor app.py
	@echo exe: weather_advisor/dist/weather_advisor.exe

clean:
	-rm -rf weather_advisor/build weather_advisor/*.spec weather_advisor/__pycache__
