.PHONY: setup run clean help

help:
	@echo Доступные команды:
	@echo   make setup  - Установить зависимости
	@echo   make run    - Запустить эксперименты
	@echo   make clean  - Удалить результаты и графики

setup:
	@echo Инициализация Go модуля...
	go mod init lab3 || echo Модуль уже инициализирован
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
	@if exist r2_comparison.png del /F r2_comparison.png
	@if exist rmse_comparison.png del /F rmse_comparison.png
	@if exist r2_vs_time.png del /F r2_vs_time.png
	@echo Очистка завершена!

