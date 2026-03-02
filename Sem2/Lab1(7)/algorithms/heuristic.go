package algorithms

import "math"

// HeuristicSearch - перебор с отсечением (LCV/Forward Checking)
func HeuristicSearch(n int) ([]int, int) {
	board := make([]int, n)
	iters := 0

	// Оценка "тесноты" позиции (число конфликтов)
	getConflicts := func(c, r int) int {
		conf := 0
		for i := 0; i < c; i++ {
			if board[i] == r || math.Abs(float64(board[i]-r)) == math.Abs(float64(i-c)) {
				conf++
			}
		}
		return conf
	}

	var solve func(col int) bool
	solve = func(col int) bool {
		if col == n {
			return true
		}

		// Сортировка строк по минимуму конфликтов (эвристика)
		type rowWeight struct {
			r, w int
		}
		rows := make([]rowWeight, n)
		for r := 0; r < n; r++ {
			rows[r] = rowWeight{r, getConflicts(col, r)}
		}

		for _, rw := range rows {
			iters++
			if rw.w == 0 { // Отсечение веток с конфликтами
				board[col] = rw.r
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
