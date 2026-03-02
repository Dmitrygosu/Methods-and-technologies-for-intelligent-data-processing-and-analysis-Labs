package algorithms

import "math/rand"

// LogicalSearch - имитация логики человека (жадный поиск с рестартами)
func LogicalSearch(n int) ([]int, int) {
	iters := 0
	maxRestarts := 1000

	for restart := 0; restart < maxRestarts; restart++ {
		board := make([]int, n)
		success := true

		for col := 0; col < n; col++ {
			found := false
			// Человек пробует ставить со случайного смещения, чтобы рестарты были разными
			startRow := rand.Intn(n)
			for i := 0; i < n; i++ {
				row := (startRow + i) % n
				iters++
				board[col] = row
				if isSafeBT(board, col) {
					found = true
					break
				}
			}

			if !found {
				success = false
				break
			}
		}

		if success {
			return board, iters
		}
	}

	return nil, iters
}
