package main

import (
	"fmt"
	"log"
	"time"

	"lab3/experiment"
	"lab3/utils"
)

func main() {
	fmt.Println("=== Лабораторная работа №3: Исследование моделей регрессии машинного обучения ===")
	fmt.Println("Начало экспериментов...")
	fmt.Println()

	startTime := time.Now()

	paramGrid := experiment.ParamGrid{
		RidgeAlpha:     []float64{0.1, 1.0, 10.0},
		LassoAlpha:     []float64{0.1, 1.0, 10.0},
		DTMaxDepth:     []int{5, 10, 15},
		RFNumTrees:     []int{10, 20, 30},
		RFMaxDepth:     []int{5, 10},
		RFMaxFeatures:  []int{2, 3},
		GBNumTrees:     []int{10, 20, 30},
		GBMaxDepth:     []int{3, 5},
		GBLearningRate: []float64{0.01, 0.1},
	}

	runner := experiment.NewExperimentRunner(paramGrid)
	results, err := runner.RunAllExperiments()
	if err != nil {
		log.Fatalf("Ошибка при выполнении экспериментов: %v", err)
	}

	err = results.SaveToJSON("results.json")
	if err != nil {
		log.Fatalf("Ошибка при сохранении результатов: %v", err)
	}

	fmt.Println()
	fmt.Printf("Эксперименты завершены за %v\n", time.Since(startTime))
	fmt.Println("Результаты сохранены в results.json")
	fmt.Println()

	fmt.Println("Генерация графиков...")

	err = utils.GenerateR2ComparisonPlot("results.json", "r2_comparison.png")
	if err != nil {
		log.Printf("Предупреждение: не удалось создать график сравнения R²: %v", err)
	} else {
		fmt.Println("r2_comparison.png создан")
	}

	err = utils.GenerateRMSEComparisonPlot("results.json", "rmse_comparison.png")
	if err != nil {
		log.Printf("Предупреждение: не удалось создать график сравнения RMSE: %v", err)
	} else {
		fmt.Println("rmse_comparison.png создан")
	}

	err = utils.GenerateR2VsTimePlot("results.json", "r2_vs_time.png")
	if err != nil {
		log.Printf("Предупреждение: не удалось создать график качества: %v", err)
	} else {
		fmt.Println("r2_vs_time.png создан")
	}

	fmt.Println()
	fmt.Println("=== Работа завершена успешно! ===")
}

