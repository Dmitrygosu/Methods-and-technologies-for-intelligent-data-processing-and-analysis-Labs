.PHONY: setup run clean help

help:
	@echo Доступные команды:
	@echo   make setup  - Установить зависимости
	@echo   make run    - Запустить эксперименты
	@echo   make clean  - Удалить результаты и графики

setup:
	@echo Инициализация Go модуля...
	go mod init lab4 || echo Модуль уже инициализирован
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
	@if exist clustering_comparison.png del /F clustering_comparison.png
	@if exist silhouette_comparison.png del /F silhouette_comparison.png
	@if exist calinski_comparison.png del /F calinski_comparison.png
	@if exist davies_comparison.png del /F davies_comparison.png
	@echo Очистка завершена!

