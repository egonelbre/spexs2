package set

import (
	"set/trie"
	"testing"
)

func TestTrieUse(t *testing.T) {
	set := trie.New()
	testUse(set, t)
}

func TestTrieMemory(t *testing.T) {
	set := trie.New()
	testMemoryUse(set, 50000, t)
}

func BenchmarkTrieAdd(b *testing.B) {
	set := trie.New()
	addValues(set, b.N)
}

func BenchmarkTrieIter(b *testing.B) {
	set := trie.New()
	b.StopTimer()
	addValues(set, b.N)
	b.StartTimer()
	iterate(set)
}
