package trie

type FullSet struct {
	Ref   *Reference
	Count int
}

func NewFullSet(ref *Reference) *FullSet {
	f := &FullSet{ref, 0}
	f.Count = len(ref.Pats)
	return f
}

func (f *FullSet) Add(idx int, pos byte) {}

func (f *FullSet) Contains(idx int, pos byte) bool {
	return idx < len(f.Ref.Pats) && int(pos) < len(f.Ref.Pats[idx].Pat)
}

func (f *FullSet) Len() int {
	return f.Count
}

func (f *FullSet) Iter() map[int]int {
	result := make(map[int]int, len(f.Ref.Pats))

	for idx, pat := range f.Ref.Pats {
		result[idx] = 2<<byte(len(pat.Pat)) - 1
	}

	return result
}
