package features

func count(arr []int, group []int) int {
	total := 0
	for _, id := range group {
		total += arr[id]
	}
	return total
}

func countf(arr []int, group []int) float64 {
	return float64(count(arr, group))
}
