package algorithms

import (
	"math"
)

// Backtracking - рекурсивный поиск с возвратом
func Backtracking(n int) ([]int, int) {
	board := make([]int, n)
	iters := 0

	var solve func(col int) bool
	solve = func(col int) bool {
		if col == n {
			return true
		}

		for row := 0; row < n; row++ {
			iters++
			board[col] = row
			if isSafeBT(board, col) {
				if solve(col + 1) {
					return true
				}
			}
		}
		return false
	}

	if solve(0) {
		return board, iters
	}
	return nil, iters
}

func isSafeBT(board []int, col int) bool {
	for i := 0; i < col; i++ {
		if board[i] == board[col] || math.Abs(float64(board[i]-board[col])) == math.Abs(float64(i-col)) {
			return false
		}
	}
	return true
}
