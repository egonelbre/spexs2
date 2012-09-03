package sets

type HashSet struct {
	data Positions
}

func NewHashSet(size int) *HashSet {
	return &HashSet{make(Positions, size)}
}

func (hs *HashSet) Add(idx int, pos int) {
	hs.data[idx] |= 1 << uint(pos)
}

func (hs *HashSet) Contains(idx int, pos int) bool {
	pv, exists := hs.data[idx]
	if !exists {
		return false
	}
	return (pv>>uint(pos))&1 == 1
}

func (hs *HashSet) Len() int {
	return len(hs.data)
}

func (hs *HashSet) Iter() Positions {
	return hs.data
}

func (hs *HashSet) Clear() {
	hs.data = nil
}

func (hs *HashSet) AddSet(g *HashSet) {
	for gidx, gval := range g.data {
		hval, exists := hs.data[gidx]
		if exists {
			hs.data[gidx] = hval | gval
		} else {
			hs.data[gidx] = gval
		}
	}
}
