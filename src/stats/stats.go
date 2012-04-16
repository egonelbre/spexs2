package stats

import "math"

func lnG(v int) float64 {
	r, _ := math.Lgamma(float64(v))
	return r
}

func gamma(v int) float64 {
	return math.Gamma(float64(v))
}

// returns probability of split of
// o - observed in set , r - observed in validation set
// O - total observed in validation, r - total items in validation set
// using logarithmic gamma function
func HypergeometricSplitLog(o int, r int, O int, R int) float64 {
	nom := lnG(O+1) + lnG(R+1) + lnG(o+r+1) + lnG(O+R-o-r+1)
	denom := lnG(o+1) + lnG(O-o+1) + lnG(r+1) + lnG(R-r+1) + lnG(O+R+1)
	return math.Exp(nom - denom)
}

// returns probability of split of
// o - observed in set , r - observed in validation set
// O - total observed in validation, r - total items in validation set
// using regular gamma function to approximate factorial
func HypergeometricSplit(o int, r int, O int, R int) float64 {
	nom := gamma(O+1) * gamma(R+1) * gamma(o+r+1) * gamma(O+R-o-r+1)
	denom := gamma(o+1) * gamma(O-o+1) * gamma(r+1) * gamma(R-r+1) * gamma(O+R+1)
	return nom / denom
}

// returns probability of split of
// choosing x items from N items
// p - probability of getting one item
func BinomialProb(x int, N int, p float64) float64 {
	nom := lnG(N + 1)
	denom := lnG(x+1) + lnG(N-x+1)
	return math.Exp(nom-denom) * math.Pow(p, float64(x)) * math.Pow(1-p, float64(N-x))
}
