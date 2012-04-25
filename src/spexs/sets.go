package spexs

type Set interface {
	Add(idx int, pos byte)
	Contains(idx int, pos byte) bool
	Len() int
	Iter() map[int]uint
}

type HashSet struct {
	data map[int]uint
}

func NewHashSet(size int) *HashSet {
	return &HashSet{make(map[int]uint, size)}
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
	return exists && (val&(1<<pos) != 0)
}

func (hs *HashSet) Len() int {
	return len(hs.data)
}

func (hs *HashSet) Iter() map[int]uint {
	return hs.data
}

func SetAddSet(h Set, g Set) {
	switch h.(type) {
	case *HashSet:
		for gidx, gval := range g.(*HashSet).data {
			hval, exists := h.(*HashSet).data[gidx]
			if exists {
				h.(*HashSet).data[gidx] = gval | hval
			} else {
				h.(*HashSet).data[gidx] = gval
			}
		}
	default:
	}
}
