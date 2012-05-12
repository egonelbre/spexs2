package binom

import (
	. "math"
	"testing"
)

type test struct {
	x, N   int
	p      float64
	result float64
}

func TestP(t *testing.T) {
	// verification result was calculated with
	// binomial(N, x) * p^x * p^(N-x)

	tests := [...]test{
		{0, 4, 0.25, 81.0 / 256.0},
		{1, 4, 0.25, 27.0 / 64.0},
		{2, 4, 0.25, 27.0 / 128.0},
	}

	for i, test := range tests {
		p := P(test.x, test.N, test.p)

		if Abs(p-test.result) > 1e-6 {
			t.Errorf("fail %v: got %v, expected %v, \nerr=%v", i, p, test.result, Abs(p-test.result))
		}
	}
}
