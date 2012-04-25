package stats

import (
	"math"
	"testing"
)

const epsilon = 0.000001

type gammaFunc func(o int, r int, O int, R int) float64

func benchGamma(b *testing.B, fn gammaFunc) {
	for v := 0; v < 1000; v += 1 {
		for r := 0; r < 1000; r += 1 {
			fn(v, r, 13000, 13000)
		}
	}
}

func BenchmarkGamma(b *testing.B) {
	benchGamma(b, HypergeometricSplit)
}

type gammaTest struct {
	o, O, r, R int
	result     float64
}

func testGamma(t *testing.T, fn gammaFunc) {
	// verification result was calculated with
	// binomial(O, o) * binomial(R, r)/binomial(O+R, o+r)

	//o, O, R, o + r
	tests := [...]gammaTest{
		{2, 41, 9, 40, 0.9969428},
		{2, 47, 9, 40, 0.9984949},
		{2, 45, 10, 30, 0.9999026},
		{1, 30, 9, 50, 0.9937611},
		{9, 45, 2, 40, 0.03885435},
	}

	for i, test := range tests {
		p := fn(test.o, test.r, test.O, test.R)

		if math.Abs(p-test.result) > epsilon {
			t.Errorf("fail %v: got %v, expected %v, \nerr=%v", i, p, test.result, math.Abs(p-test.result))
		}
	}
}

func TestGamma(t *testing.T) {
	testGamma(t, HypergeometricSplit)
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

		if math.Abs(p-test.result) > epsilon {
			t.Errorf("fail %v: got %v, expected %v, \nerr=%v", i, p, test.result, math.Abs(p-test.result))
		}
	}
}
