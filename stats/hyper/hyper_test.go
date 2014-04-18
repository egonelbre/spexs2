package hyper

import (
	"math"
	"testing"
)

type splitFunc func(o int, r int, O int, R int) float64

func benchHyper(b *testing.B, fn splitFunc) {
	for i := 0; i < b.N; i += 1 {
		for v := 0; v < 1000; v += 13 {
			for r := 0; r < 1000; r += 15 {
				fn(v, r, 13000, 13000)
			}
		}
	}
}

func BenchmarkComplementCdf(b *testing.B) {
	benchHyper(b, ComplementCdf)
}

func BenchmarkComplementCdfApprox(b *testing.B) {
	benchHyper(b, ComplementCdfApprox)
}

func BenchmarkComplementCdfSlow(b *testing.B) {
	benchHyper(b, ComplementCdfSlow)
}

type test struct {
	chosenA, totalA, chosenB, totalB int

	expected float64
}

func testHyper(t *testing.T, fn splitFunc, epsilon float64) {
	// verification result was calculated with R
	// phyper(o-1, O, R, o+r, lower.tail = F, log.p = F)
	tests := [...]test{
		{2, 41, 9, 40, 0.9969428},
		{2, 47, 9, 40, 0.9984949},
		{2, 45, 10, 30, 0.9999026},
		{1, 30, 9, 50, 0.9937611},
		{9, 45, 2, 40, 0.03885435},
		{1700, 1000000, 5000, 3000000, 0.244143},
		{1700, 1000000, 3000, 3000000, 6.340578e-65},
	}

	for i, test := range tests {
		p := fn(test.chosenA, test.chosenB, test.totalA, test.totalB)

		diff := math.Abs(p/test.expected - 1)
		if diff > epsilon {
			t.Errorf("fail %v: got %v, expected %v, \nerr=%v", i, p, test.expected, diff)
		}
	}
}

func TestComplementCdf(t *testing.T) {
	testHyper(t, ComplementCdf, 1e-6)
}

func TestComplementCdfSlow(t *testing.T) {
	testHyper(t, ComplementCdfSlow, 1e-6)
}

func TestComplementCdfApprox(t *testing.T) {
	testHyper(t, ComplementCdfApprox, 1e-5)
}
