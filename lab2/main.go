package main

import (
	"fmt"
	"log"
	"time"

	"lab2/experiment"
	"lab2/utils"
)

func main() {
	fmt.Println("=== Лабораторная работа №2: Исследование методов классификации и кластеризации ===")
	fmt.Println("Начало экспериментов...")
	fmt.Println()

	startTime := time.Now()

	paramGrid := experiment.ParamGrid{
		KNNK:           []int{3, 5, 7},
		DTMaxDepth:     []int{5, 10, 15},
		LRLearningRate: []float64{0.01, 0.1},
		LRIterations:   []int{100, 500},
		KMeansK:        []int{2, 3, 4},
		KMeansMaxIter:  []int{100, 200},
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

	err = utils.GenerateClassificationComparisonPlot("results.json", "classification_comparison.png")
	if err != nil {
		log.Printf("Предупреждение: не удалось создать график сравнения классификаторов: %v", err)
	} else {
		fmt.Println("classification_comparison.png создан")
	}

	err = utils.GenerateClusteringComparisonPlot("results.json", "clustering_comparison.png")
	if err != nil {
		log.Printf("Предупреждение: не удалось создать график сравнения кластеризации: %v", err)
	} else {
		fmt.Println("clustering_comparison.png создан")
	}

	err = utils.GenerateAccuracyVsTimePlot("results.json", "accuracy_vs_time.png")
	if err != nil {
		log.Printf("Предупреждение: не удалось создать график точности: %v", err)
	} else {
		fmt.Println("accuracy_vs_time.png создан")
	}

	fmt.Println("Генерация матриц ошибок для лучших моделей...")

	err = utils.GenerateConfusionMatrixPlot("results.json", "confusion_matrix_naive_bayes.png", "naive_bayes")
	if err != nil {
		log.Printf("Предупреждение: не удалось создать матрицу ошибок для Naive Bayes: %v", err)
	} else {
		fmt.Println("confusion_matrix_naive_bayes.png создан")
	}

	err = utils.GenerateConfusionMatrixPlot("results.json", "confusion_matrix_knn.png", "knn")
	if err != nil {
		log.Printf("Предупреждение: не удалось создать матрицу ошибок для KNN: %v", err)
	} else {
		fmt.Println("confusion_matrix_knn.png создан")
	}

	err = utils.GenerateConfusionMatrixPlot("results.json", "confusion_matrix_decision_tree.png", "decision_tree")
	if err != nil {
		log.Printf("Предупреждение: не удалось создать матрицу ошибок для Decision Tree: %v", err)
	} else {
		fmt.Println("confusion_matrix_decision_tree.png создан")
	}

	err = utils.GenerateConfusionMatrixPlot("results.json", "confusion_matrix_logistic_regression.png", "logistic_regression")
	if err != nil {
		log.Printf("Предупреждение: не удалось создать матрицу ошибок для Logistic Regression: %v", err)
	} else {
		fmt.Println("confusion_matrix_logistic_regression.png создан")
	}

	fmt.Println()
	fmt.Println("=== Работа завершена успешно! ===")
}
