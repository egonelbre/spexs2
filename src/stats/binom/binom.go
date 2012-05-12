package binom

import (
	. "math"
)

func lnG(v int) float64 {
	r, _ := Lgamma(float64(v))
	return r
}

func gamma(v int) float64 {
	return Gamma(float64(v))
}

// returns probability of split of
// choosing x items from N items
// p - probability of getting one item
func P(x int, N int, p float64) float64 {
	nom := lnG(N + 1)
	denom := lnG(x+1) + lnG(N-x+1)
	return Exp(nom-denom) * Pow(p, float64(x)) * Pow(1-p, float64(N-x))
}
