package algorithms

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

type Msg struct {
	AgentID int
	Row     int
}

// MultiAgentSearch - MAS на горутинах с эвристикой Min-Conflicts и защитой от дедлока
func MultiAgentSearch(n int, timeout time.Duration) ([]int, int) {
	rand.Seed(time.Now().UnixNano())
	board := make([]int, n)
	for i := 0; i < n; i++ {
		board[i] = rand.Intn(n)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // Гарантирует отмену контекста при выходе

	var mu sync.Mutex
	iters := 0
	comm := make(chan Msg, n*10)

	// Запуск агентов
	for i := 0; i < n; i++ {
		go func(id int) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					mu.Lock()
					conf := countConflicts(board, id)
					var changed bool
					var newRow int

					if conf > 0 {
						minConf := n + 1
						var bestRows []int

						for r := 0; r < n; r++ {
							if r == board[id] {
								continue
							}
							oldRow := board[id]
							board[id] = r
							c := countConflicts(board, id)
							if c < minConf {
								minConf = c
								bestRows = []int{r}
							} else if c == minConf {
								bestRows = append(bestRows, r)
							}
							board[id] = oldRow
						}

						if rand.Float32() < 0.10 {
							board[id] = rand.Intn(n)
						} else if len(bestRows) > 0 {
							board[id] = bestRows[rand.Intn(len(bestRows))]
						}

						newRow = board[id]
						changed = true
					}
					mu.Unlock()

					if changed {
						// Критически важно: проверяем контекст при отправке,
						// чтобы не зависнуть, если координатор уже ушел.
						select {
						case comm <- Msg{AgentID: id, Row: newRow}:
						case <-ctx.Done():
							return
						default:
							// Канал забит, пропускаем отправку, чтобы избежать дедлока
						}
					}
					time.Sleep(10 * time.Microsecond)
				}
			}
		}(i)
	}

	// Координатор
	for {
		select {
		case <-ctx.Done():
			return nil, iters
		case <-comm:
			iters++
			mu.Lock()
			totalConf := 0
			for i := 0; i < n; i++ {
				totalConf += countConflicts(board, i)
			}
			if totalConf == 0 {
				res := make([]int, n)
				copy(res, board)
				mu.Unlock()
				// После return сработает defer cancel(), и все агенты завершатся
				return res, iters
			}
			mu.Unlock()
		}
	}
}

func countConflicts(board []int, col int) int {
	c := 0
	for i := 0; i < len(board); i++ {
		if i == col {
			continue
		}
		if board[i] == board[col] {
			c++
		}
		diffX := col - i
		if diffX < 0 {
			diffX = -diffX
		}
		diffY := board[col] - board[i]
		if diffY < 0 {
			diffY = -diffY
		}
		if diffX == diffY {
			c++
		}
	}
	return c
}
