package spexs

type Set interface {
	Add(idx int, pos byte)
	Contains(idx int, pos byte) bool
	Len() int
	Iter() map[int]int
}

type HashSet struct {
	data map[int]int
}

func NewHashSet(size int) *HashSet {
	return &HashSet{make(map[int]int, size)}
}

func (hs *HashSet) Add(idx int, pos byte) {
	val, exists := hs.data[idx]
	if !exists {
		val = 0
	}
	hs.data[idx] = val | (1 << pos)
}

func (hs *HashSet) Contains(idx int, pos byte) bool {
	val, exists := hs.data[idx]
	return exists && (val & (1 << pos) != 0)
}

func (hs *HashSet) Len() int {
	return len(hs.data)
}

func (hs *HashSet) Iter() map[int]int {
	return hs.data;
}

func SetAddSet(h Set, g Set) {
	switch h.(type){
	case *HashSet : 
		for gidx, gval := range g.(*HashSet).data {
			hval, exists := h.(*HashSet).data[gidx]
			if exists {
				h.(*HashSet).data[gidx] = gval | hval
			} else {
				h.(*HashSet).data[gidx] = gval
			}
		}
		default :
	}
}

type FullSet struct {
	Ref   *UnicodeReference
	Count int
}

func NewFullSet(ref *UnicodeReference) *FullSet {
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
		result[idx] = 2 << byte(len(pat.Pat)) - 1
	}

	return result
}

