package utils

import (
	"math"
	"testing"
)

type gammaFunc func(o int, r int, O int, R int) float64

func benchHyper(b *testing.B, fn gammaFunc) {
	for v := 0; v < 100; v += 1 {
		for r := 0; r < 100; r += 1 {
			fn(v, r, 13000, 13000)
		}
	}
}

func BenchmarkHyper(b *testing.B) {
	benchHyper(b, HypergeometricSplit)
}

func BenchmarkHyperApprox(b *testing.B) {
	benchHyper(b, HypergeometricSplitApprox)
}

func BenchmarkHyperSlow(b *testing.B) {
	benchHyper(b, HypergeometricSplitSlow)
}

type gammaTest struct {
	o, O, r, R int
	result     float64
}

func testHyper(t *testing.T, fn gammaFunc, epsilon float64) {
	// verification result was calculated with R
	// phyper(o-1, O, R, o+r, lower.tail = F, log.p = F)
	tests := [...]gammaTest{
		{2, 41, 9, 40, 0.9969428},
		{2, 47, 9, 40, 0.9984949},
		{2, 45, 10, 30, 0.9999026},
		{1, 30, 9, 50, 0.9937611},
		{9, 45, 2, 40, 0.03885435},
		{1700, 1000000, 5000, 3000000, 0.244143},
		{1700, 1000000, 3000, 3000000, 6.340578e-65},
	}

	for i, test := range tests {
		p := fn(test.o, test.r, test.O, test.R)

		if math.Abs(p-test.result) > epsilon {
			t.Errorf("fail %v: got %v, expected %v, \nerr=%v", i, p, test.result, math.Abs(p-test.result))
		}
	}
}

func TestHypergeometricSplit(t *testing.T) {
	testHyper(t, HypergeometricSplit, 1e-6)
}

func TestHypergeometricSplitSlow(t *testing.T) {
	testHyper(t, HypergeometricSplitSlow, 1e-6)
}

func TestHypergeometricSplitApprox(t *testing.T) {
	testHyper(t, HypergeometricSplitApprox, 1e-5)
}

type binomTest struct {
	x, N   int
	p      float64
	result float64
}

func TestBinomial(t *testing.T) {
	// verification result was calculated with
	// binomial(N, x) * p^x * p^(N-x)

	tests := [...]binomTest{
		{0, 4, 0.25, 81.0 / 256.0},
		{1, 4, 0.25, 27.0 / 64.0},
		{2, 4, 0.25, 27.0 / 128.0},
	}

	for i, test := range tests {
		p := BinomialProb(test.x, test.N, test.p)

		if math.Abs(p-test.result) > 1e-6 {
			t.Errorf("fail %v: got %v, expected %v, \nerr=%v", i, p, test.result, math.Abs(p-test.result))
		}
	}
}
