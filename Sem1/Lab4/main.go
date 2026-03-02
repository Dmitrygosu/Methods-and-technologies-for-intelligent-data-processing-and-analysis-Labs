package main

import (
	"fmt"
	"log"
	"time"

	"lab4/experiment"
	"lab4/utils"
)

func main() {
	fmt.Println("=== Лабораторная работа №4: Исследование методов кластеризации машинного обучения ===")
	fmt.Println("Начало экспериментов...")
	fmt.Println()

	startTime := time.Now()

	paramGrid := experiment.ParamGrid{
		KMeansK:          []int{2, 3, 4, 5},
		KMeansMaxIter:    []int{100, 200},
		DBSCANEps:        []float64{0.3, 0.5, 0.7},
		DBSCANMinSamples: []int{3, 5},
		HierarchicalK:    []int{2, 3, 4},
		GMMComponents:    []int{2, 3, 4},
		GMMMaxIter:       []int{50, 100},
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

	err = utils.GenerateClusteringComparisonPlot("results.json", "clustering_comparison.png")
	if err != nil {
		log.Printf("Предупреждение: не удалось создать график сравнения кластеризации: %v", err)
	} else {
		fmt.Println("clustering_comparison.png создан")
	}

	err = utils.GenerateSilhouetteComparisonPlot("results.json", "silhouette_comparison.png")
	if err != nil {
		log.Printf("Предупреждение: не удалось создать график сравнения индекса силуэта: %v", err)
	} else {
		fmt.Println("silhouette_comparison.png создан")
	}

	err = utils.GenerateCalinskiComparisonPlot("results.json", "calinski_comparison.png")
	if err != nil {
		log.Printf("Предупреждение: не удалось создать график сравнения индекса Калинского-Харабаша: %v", err)
	} else {
		fmt.Println("calinski_comparison.png создан")
	}

	err = utils.GenerateDaviesComparisonPlot("results.json", "davies_comparison.png")
	if err != nil {
		log.Printf("Предупреждение: не удалось создать график сравнения индекса Дэвиса-Болдуина: %v", err)
	} else {
		fmt.Println("davies_comparison.png создан")
	}

	fmt.Println()
	fmt.Println("=== Работа завершена успешно! ===")
}
