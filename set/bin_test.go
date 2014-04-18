package set

import (
	"set/bin"
	"testing"
)

func TestBinUse(t *testing.T) {
	set := bin.New()
	testUse(set, t)
}

func TestBinMemory(t *testing.T) {
	set := bin.New()
	testMemoryUse(set, 50000, t)
}

func BenchmarkBinAdd(b *testing.B) {
	set := bin.New()
	addValues(set, b.N)
}

func BenchmarkBinIter(b *testing.B) {
	set := bin.New()
	b.StopTimer()
	addValues(set, b.N)
	b.StartTimer()
	iterate(set)
}
