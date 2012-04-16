package stats

import "testing"

func BenchmarkLogGamma(b *testing.B) {
	for v := 0; v < 1000; v += 1 {
		for r := 0; r < 1000; r += 1 {
			HypergeometricSplitLog(v, r, 13000, 13000)
		}
	}
}

func BenchmarkGamma(b *testing.B) {
	for v := 0; v < 1000; v += 1 {
		for r := 0; r < 1000; r += 1 {
			HypergeometricSplit(v, r, 13000, 13000)
		}
	}
}
