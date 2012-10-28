package features

func countOnly(count []int, group []int) float64 {
	var total float64 = 0.0
	for _, id := range group {
		total += count[id]
	}
	return total
}
