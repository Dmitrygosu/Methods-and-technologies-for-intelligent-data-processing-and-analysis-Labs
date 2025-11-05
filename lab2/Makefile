.PHONY: setup run clean help

help:
	@echo Доступные команды:
	@echo   make setup  - Установить зависимости
	@echo   make run    - Запустить эксперименты
	@echo   make clean  - Удалить результаты и графики

setup:
	@echo Инициализация Go модуля...
	go mod init lab2 || echo Модуль уже инициализирован
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
	@if exist classification_comparison.png del /F classification_comparison.png
	@if exist clustering_comparison.png del /F clustering_comparison.png
	@if exist accuracy_vs_time.png del /F accuracy_vs_time.png
	@if exist confusion_matrix_*.png del /F confusion_matrix_*.png
	@echo Очистка завершена!

