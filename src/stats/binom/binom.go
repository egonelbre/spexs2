package binom

import (
	stat "code.google.com/p/go-fn/fn"
)

// the probability
// returns probability of
// choosing at least k items from n items
// p - probability of getting one item
func P(k, n int, p float64) float64 {
	return stat.BetaIncReg(float64(k+1), float64(n-k), p)
}
