package hyper

import (
	"math"
	"testing"
)

type splitFunc func(o int, r int, O int, R int) float64

func benchHyper(b *testing.B, fn splitFunc) {
	for v := 0; v < 100; v += 1 {
		for r := 0; r < 100; r += 1 {
			fn(v, r, 13000, 13000)
		}
	}
}

func BenchmarkSplit(b *testing.B) {
	benchHyper(b, Split)
}

func BenchmarkSplitApprox(b *testing.B) {
	benchHyper(b, SplitApprox)
}

func BenchmarkSplitSlow(b *testing.B) {
	benchHyper(b, SplitSlow)
}

type test struct {
	o, O, r, R int
	result     float64
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
		p := fn(test.o, test.r, test.O, test.R)

		diff := math.Abs(p / test.result - 1)
		if diff > epsilon {
			t.Errorf("fail %v: got %v, expected %v, \nerr=%v", i, p, test.result, diff)
		}
	}
}

func TestSplit(t *testing.T) {
	testHyper(t, Split, 1e-6)
}

func TestSplitSlow(t *testing.T) {
	testHyper(t, SplitSlow, 1e-6)
}

func TestSplitApprox(t *testing.T) {
	testHyper(t, SplitApprox, 1e-5)
}
