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

func BenchmarkLogGamma(b *testing.B) {
	benchGamma(b, HypergeometricSplitLog)
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

	tests := [...]gammaTest{
		{2, 5, 2, 10, 30.0 / 91.0},
		{2, 45, 9, 40, 71632.0 / 5645577.0},
	}

	for i, test := range tests {
		p := fn(test.o, test.r, test.O, test.R)
		p2 := fn(test.r, test.o, test.R, test.O) // since the test must be symmetric

		if math.Abs(p-test.result) > epsilon || math.Abs(p2-test.result) > epsilon {
			t.Errorf("fail %v: got (%v, %v), expected %v, \nerr=(%v,%v)", i, p, p2, test.result, math.Abs(p-test.result), math.Abs(p2-test.result))
		}
	}
}

func TestLogGamma(t *testing.T) {
	testGamma(t, HypergeometricSplitLog)
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
