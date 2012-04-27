package trie

import (
	. "spexs"
	"math/big"
)

type FullSet struct {
	Ref   *Reference
	Count int
}

func NewFullSet(ref *Reference) *FullSet {
	f := &FullSet{ref, 0}
	f.Count = len(ref.Seqs)
	return f
}

func (f *FullSet) Add(idx int, pos int) {}

func (f *FullSet) Contains(idx int, pos int) bool {
	return idx < len(f.Ref.Seqs) && int(pos) < len(f.Ref.Seqs[idx].Pat)
}

func (f *FullSet) Len() int {
	return f.Count
}

func (hs *FullSet) Clear() {
}

func (f *FullSet) Iter() Positions {
	result := make(Positions, len(f.Ref.Seqs))

	for idx, pat := range f.Ref.Seqs {
		i := big.NewInt(2)
		i.Lsh(i, uint(len(pat.Pat)))
		i.Sub(i, big.NewInt(1))
		result[idx] = i
		// 2<<byte(len(pat.Pat)) - 1
	}

	return result
}
