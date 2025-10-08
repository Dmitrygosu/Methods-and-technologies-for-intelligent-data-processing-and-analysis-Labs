package main

import (
	"fmt"
	"log"
	"time"

	"lab1/experiment"
	"lab1/utils"
)

func main() {
	fmt.Println("=== Лабораторная работа №1: Исследование генетического алгоритма ===")
	fmt.Println("Начало экспериментов...")
	fmt.Println()

	startTime := time.Now()

	paramGrid := experiment.ParamGrid{
		PopulationSizes: []int{50, 100, 200},
		MaxGenerations:  []int{25, 50, 75},
		CrossoverProbs:  []float64{0.6, 0.8},
		MutationProbs:   []float64{0.01, 0.05, 0.1},
		CrossoverTypes:  []string{"onepoint", "uniform"},
		ElitismCounts:   []int{2, 5},
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

	err = utils.GenerateTimeComparisonPlot("results.json", "time_comparison.png")
	if err != nil {
		log.Printf("Предупреждение: не удалось создать график времени: %v", err)
	} else {
		fmt.Println("time_comparison.png создан")
	}

	err = utils.GenerateConvergencePlot("results.json", "convergence_array.png")
	if err != nil {
		log.Printf("Предупреждение: не удалось создать график сходимости: %v", err)
	} else {
		fmt.Println("convergence_array.png создан")
	}

	err = utils.GenerateAccuracyVsTimePlot("results.json", "accuracy_vs_time.png")
	if err != nil {
		log.Printf("Предупреждение: не удалось создать график точности: %v", err)
	} else {
		fmt.Println("accuracy_vs_time.png создан")
	}

	err = utils.GenerateEfficiencyComparisonPlot("results.json", "efficiency_comparison.png")
	if err != nil {
		log.Printf("Предупреждение: не удалось создать график эффективности: %v", err)
	} else {
		fmt.Println("efficiency_comparison.png создан")
	}

	fmt.Println()
	fmt.Println("=== Работа завершена успешно! ===")
}
