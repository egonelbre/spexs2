package spexs

type Set interface {
	Add(idx int, pos byte)
	Contains(idx int, pos byte) bool
	Length() int
	Iter() (chan int, chan int)
}

type HashSet struct {
	data map[int]int
}

func NewHashSet() *HashSet {
	return &HashSet{make(map[int]int)}
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

func (hs *HashSet) Length() int {
	return len(hs.data)
}

func (hs *HashSet) Iter() (chan int, chan int) {
	indices := make(chan int, 100)
	poss := make(chan int, 100)

	go func(){
		for idx, pos := range hs.data {
			indices <- idx
			poss <- pos
		}
		close(indices)
		close(poss)
	}()
	return indices, poss
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
	for _, p := range ref.Pats {
		f.Count += p.Count
	}
	return f
}

func (f *FullSet) Add(idx int, pos byte) {}

func (f *FullSet) Contains(idx int, pos byte) bool {
	return idx < len(f.Ref.Pats) && int(pos) < len(f.Ref.Pats[idx].Pat)
}

func (f *FullSet) Length() int {
	return f.Count
}

func (f *FullSet) Iter() (chan int, chan int) {
	indices := make(chan int, 100)
	poss := make(chan int, 100)

	go func(){
		for idx, pat := range f.Ref.Pats {
			indices <- idx
			poss <- (2 << byte(len(pat.Pat) - 1))
		}
		close(indices)
		close(poss)
	}()
	return indices, poss
}

