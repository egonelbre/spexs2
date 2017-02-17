package binom

import (
	"github.com/ematvey/go-fn/fn"
)

// the probability
// returns probability of
// gettting at least "successes" from "total" items
// p - probability of getting one success
func ComplementCdf(successes, total int, p float64) float64 {
	return fn.BetaIncReg(float64(successes+1), float64(total-successes), p)
}
