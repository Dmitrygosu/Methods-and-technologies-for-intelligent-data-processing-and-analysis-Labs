package main

import (
	"encoding/json"
	"fmt"
	"os"
	"queens-lab/algorithms"
	"time"
)

func main() {
	ns := []int{8, 10, 12}
	var allResults []Result

	fmt.Println("--- Запуск бенчмарков для задачи о ферзях ---")

	for _, n := range ns {
		fmt.Printf("\nРазмер доски: %dx%d\n", n, n)

		// 1. Brute Force (с лимитом, т.к. n! растет быстро)
		resBF := RunBenchmark("BruteForce (Naive)", n, func(ni int) ([]int, int) {
			return algorithms.BruteForceNaive(ni, 1000000)
		})
		allResults = append(allResults, resBF)
		fmt.Printf("BF:  Time: %.3f ms, Iter: %d\n", resBF.Time, resBF.Iter)

		// 2. Backtracking
		resBT := RunBenchmark("Backtracking", n, algorithms.Backtracking)
		allResults = append(allResults, resBT)
		fmt.Printf("BT:  Time: %.3f ms, Iter: %d\n", resBT.Time, resBT.Iter)

		// 3. Heuristic
		resH := RunBenchmark("Heuristic (Pruning)", n, algorithms.HeuristicSearch)
		allResults = append(allResults, resH)
		fmt.Printf("H:   Time: %.3f ms, Iter: %d\n", resH.Time, resH.Iter)

		// 4. Logical
		resL := RunBenchmark("Logical (Greedy)", n, algorithms.LogicalSearch)
		allResults = append(allResults, resL)
		fmt.Printf("L:   Time: %.3f ms, Iter: %d\n", resL.Time, resL.Iter)

		// 5. MAS
		resMAS := RunBenchmark("Multi-Agent", n, func(ni int) ([]int, int) {
			return algorithms.MultiAgentSearch(ni, 2*time.Second)
		})
		allResults = append(allResults, resMAS)
		fmt.Printf("MAS: Time: %.3f ms, Iter: %d\n", resMAS.Time, resMAS.Iter)
	}

	// Сохранение в JSON для визуализации
	file, _ := json.MarshalIndent(allResults, "", "  ")
	_ = os.WriteFile("results.json", file, 0644)
	fmt.Println("\nРезультаты сохранены в results.json")
}
