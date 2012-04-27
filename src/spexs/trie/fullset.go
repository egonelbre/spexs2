package trie

type FullSet struct {
	Ref   *Reference
	Count int
}

func NewFullSet(ref *Reference) *FullSet {
	f := &FullSet{ref, 0}
	f.Count = len(ref.Seqs)
	return f
}

func (f *FullSet) Add(idx int, pos byte) {}

func (f *FullSet) Contains(idx int, pos byte) bool {
	return idx < len(f.Ref.Seqs) && int(pos) < len(f.Ref.Seqs[idx].Pat)
}

func (f *FullSet) Len() int {
	return f.Count
}

func (f *FullSet) Iter() map[int]uint64 {
	result := make(map[int]uint64, len(f.Ref.Seqs))

	for idx, pat := range f.Ref.Seqs {
		result[idx] = 2<<byte(len(pat.Pat)) - 1
	}

	return result
}
