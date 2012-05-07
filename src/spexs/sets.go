package spexs

type BitVector uint16
type Positions map[int]BitVector

type Set interface {
	Add(idx int, pos uint)
	Contains(idx int, pos uint) bool
	Len() int
	Iter() Positions
	Clear()
}

type HashSet struct {
	data Positions
}

func NewHashSet(size int) *HashSet {
	return &HashSet{make(Positions, size)}
}

func (hs *HashSet) Add(idx int, pos uint) {
	val, exists := hs.data[idx]
	if !exists {
		val = 0
		hs.data[idx] = 0
	}
	val |= 1 << pos
	hs.data[idx] = val
}

func (hs *HashSet) Contains(idx int, pos uint) bool {
	val, exists := hs.data[idx]
	return exists && (val & (1 << pos) != BitVector(0))
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

func (hs *HashSet) AddSet(g HashSet) {
	for gidx, gval := range g.data {
		hval, exists := hs.data[gidx]
		if exists {
			hs.data[gidx] = hval | gval
		} else {
			hs.data[gidx] = gval
		}
	}
}
