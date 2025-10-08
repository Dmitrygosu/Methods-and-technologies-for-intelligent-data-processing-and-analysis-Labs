# На Windows для запуска нужно установить инструмент make через choco install make 
.PHONY: setup run clean help

help:
	@echo Доступные команды:
	@echo   make setup  - Установить зависимости
	@echo   make run    - Запустить эксперименты
	@echo   make clean  - Удалить результаты и графики

setup:
	@echo Инициализация Go модуля...
	go mod init lab1 || echo Модуль уже инициализирован
	@echo Установка зависимостей...
	go get gonum.org/v1/plot
	go mod tidy
	@echo Зависимости установлены успешно!

run:
	@echo Запуск экспериментов...
	go run main.go
	@echo Готово! Проверьте results.json и графики (.png)

clean:
	@echo Очистка результатов...
	@if exist results.json del /F results.json
	@if exist time_comparison.png del /F time_comparison.png
	@if exist convergence_array.png del /F convergence_array.png
	@if exist accuracy_vs_time.png del /F accuracy_vs_time.png
	@if exist efficiency_comparison.png del /F efficiency_comparison.png
	@echo Очистка завершена!

