package utils

import "math"

func lnG(v int) float64 {
	r, _ := math.Lgamma(float64(v))
	return r
}

func gamma(v int) float64 {
	return math.Gamma(float64(v))
}

// returns probability of split of
// o - observed in input , r - observed in validation set
// O - total items in input, R - total items in validation set
// using logarithmic gamma function
// TODO: limits test
func HypergeometricSplit(o int, r int, O int, R int) float64 {
	total := 0.0
	lSOR := lnG(O+1) + lnG(R+1)
	lOR := lnG(O + R + 1)
	for r >= 0 {
		nom := lSOR + lnG(o+r+1) + lnG(O+R-o-r+1)
		denom := lnG(o+1) + lnG(O-o+1) + lnG(r+1) + lnG(R-r+1) + lOR
		add := math.Exp(nom - denom)
		total += add
		r -= 1
		o += 1
	}
	return total
}

// returns probability of split of
// o - observed in input , r - observed in validation set
// O - total items in input, R - total items in validation set
// using logarithmic gamma function
// TODO: limits test
func HypergeometricSplitApprox(o int, r int, O int, R int) float64 {
	total := 0.0
	lSOR := lnG(O+1) + lnG(R+1)
	lOR := lnG(O + R + 1)
	for r >= 0 {
		nom := lSOR + lnG(o+r+1) + lnG(O+R-o-r+1)
		denom := lnG(o+1) + lnG(O-o+1) + lnG(r+1) + lnG(R-r+1) + lOR
		add := math.Exp(nom - denom)
		total += add
		// add will first increase then decrease
		// this means add must go through the maximum point of the distribution
		// if add < total/1eN we must be already decreasing
		// and therefore total/1eN means approximately N digits precision
		if add < total/1e16 {
			break
		}
		r -= 1
		o += 1
	}
	return total
}

// returns probability of split of
// o - observed in input , r - observed in validation set
// O - total items in input, R - total items in validation set
// using logarithmic gamma function
func HypergeometricSplitDown(o int, r int, O int, R int) float64 {
	total := 0.0
	lSOR := lnG(O+1) + lnG(R+1)
	lOR := lnG(O + R + 1)
	for o >= 0 {
		nom := lSOR + lnG(o+r+1) + lnG(O+R-o-r+1)
		denom := lnG(o+1) + lnG(O-o+1) + lnG(r+1) + lnG(R-r+1) + lOR
		total += math.Exp(nom - denom)

		r += 1
		o -= 1
	}
	return total
}

// returns probability of split of
// choosing x items from N items
// p - probability of getting one item
func BinomialProb(x int, N int, p float64) float64 {
	nom := lnG(N + 1)
	denom := lnG(x+1) + lnG(N-x+1)
	return math.Exp(nom-denom) * math.Pow(p, float64(x)) * math.Pow(1-p, float64(N-x))
}
