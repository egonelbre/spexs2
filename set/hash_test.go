package set

import (
	"set/hash"
	"testing"
)

func TestHashUse(t *testing.T) {
	set := hash.New()
	testUse(set, t)
}

func TestHashMemory(t *testing.T) {
	set := hash.New()
	testMemoryUse(set, 50000, t)
}

func BenchmarkHashAdd(b *testing.B) {
	set := hash.New()
	addValues(set, b.N)
}

func BenchmarkHashIter(b *testing.B) {
	set := hash.New()
	b.StopTimer()
	addValues(set, b.N)
	b.StartTimer()
	iterate(set)
}
