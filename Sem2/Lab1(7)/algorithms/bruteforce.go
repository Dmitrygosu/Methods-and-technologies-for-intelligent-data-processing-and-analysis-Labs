package algorithms

// BruteForceNaive - полный перебор всех перестановок.
// n - размер доски, maxIter - лимит итераций
func BruteForceNaive(n int, maxIter int) ([]int, int) {
	iters := 0
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = i
	}

	result := []int{}

	var generate func(k int) bool
	generate = func(k int) bool {
		if k == n {
			iters++
			if isAllSafe(nums) {
				result = make([]int, n)
				copy(result, nums)
				return true
			}
			return iters > maxIter
		}

		for i := k; i < n; i++ {
			if iters > maxIter {
				return true
			}
			nums[k], nums[i] = nums[i], nums[k]
			if generate(k + 1) {
				return true
			}
			nums[k], nums[i] = nums[i], nums[k]
		}
		return false
	}

	generate(0)
	if len(result) > 0 {
		return result, iters
	}
	return nil, iters
}

// Полная проверка доски
func isAllSafe(board []int) bool {
	n := len(board)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			// Проверка диагоналей (горизонтали исключены перестановками)
			diffX := j - i
			diffY := board[j] - board[i]
			if diffY < 0 {
				diffY = -diffY
			}
			if diffX == diffY {
				return false
			}
		}
	}
	return true
}
