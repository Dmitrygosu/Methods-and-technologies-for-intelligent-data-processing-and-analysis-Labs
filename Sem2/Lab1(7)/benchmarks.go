package main

import (
	"fmt"
	"time"
)

type Result struct {
	Algo    string  `json:"algo"`
	N       int     `json:"n"`
	Time    float64 `json:"time_ms"`
	Iter    int     `json:"iterations"`
	Success bool    `json:"success"`
}

// RunBenchmark - запуск замера для одного алгоритма
func RunBenchmark(name string, n int, f func(int) ([]int, int)) Result {
	start := time.Now()
	board, iters := f(n)
	elapsed := time.Since(start)

	success := board != nil
	return Result{
		Algo:    name,
		N:       n,
		Time:    float64(elapsed.Microseconds()) / 1000.0,
		Iter:    iters,
		Success: success,
	}
}

// Печать доски (для отладки/демо)
func PrintBoard(board []int) {
	if board == nil {
		fmt.Println("Решение не найдено")
		return
	}
	n := len(board)
	for r := 0; r < n; r++ {
		for c := 0; c < n; c++ {
			if board[c] == r {
				fmt.Print("Q ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
}
