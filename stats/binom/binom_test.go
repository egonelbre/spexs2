package binom

import (
	"math"
	"testing"
)

type test struct {
	k, n   int
	p      float64
	result float64
}

func BenchmarkBinom(b *testing.B) {
	for i := 0; i < b.N; i ++ {
		for v := 0; v < 1000; v ++3 {
			for r := 0; r < 1000; r ++5 {
				ComplementCdf(v, 13000, float64(r)/float64(13000))
			}
		}
	}
}

func TestComplementCdf(t *testing.T) {
	// verification result was calculated with
	// pbinom(k, N, p, lower.tail=F, log.P = F)

	tests := [...]test{
		{0, 4, 0.25, 0.6835937},
		{5, 100, 0.01, 0.0005345345},
		{10, 100, 0.01, 6.25552e-9},
		{39, 98, 23.0 / 1000.0, 3.883636e-39},
	}

	for _, test := range tests {
		p := ComplementCdf(test.k, test.n, test.p)
		if math.IsNaN(p) {
			t.Errorf("got NaN from: k=%v N=%v p=%v", test.k, test.n, test.p)
		}

		diff := math.Abs(1 - p/test.result)
		if diff > 1e-5 {
			t.Errorf("failed k=%v N=%v p=%v: got %v, expected %v, \nerr=%v", test.k, test.n, test.p, p, test.result, diff)
		}
	}
}
