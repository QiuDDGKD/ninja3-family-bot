package tools

import "math/rand/v2"

// 随机 n 个小于 max 且不重复的整数
func RandInts(n, max int) []int {
	if n > max {
		n = max
	}

	result := make([]int, 0, n)
	seen := make(map[int]struct{})

	for len(result) < n {
		num := rand.Int() % max
		if _, exists := seen[num]; !exists {
			seen[num] = struct{}{}
			result = append(result, num)
		}
	}

	return result
}
