package spexs

import "math/big"

type Positions map[int]*big.Int

type Set interface {
	Add(idx int, pos int)
	Contains(idx int, pos int) bool
	Len() int
	Iter() Positions
	Clear()
}

type HashSet struct {
	data map[int]*big.Int
}

func NewHashSet(size int) *HashSet {
	return &HashSet{make(Positions, size)}
}

func (hs *HashSet) Add(idx int, pos int) {
	val, exists := hs.data[idx]
	if !exists {
		val = big.NewInt(0)
		hs.data[idx] = val
	}
	val.SetBit(val, pos, 1)
}

func (hs *HashSet) Contains(idx int, pos int) bool {
	val, exists := hs.data[idx]
	return exists && (val.Bit(pos) > 0)
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

func SetAddSet(h Set, g Set) {
	switch h.(type) {
	case *HashSet:
		for gidx, gval := range g.(*HashSet).data {
			hval, exists := h.(*HashSet).data[gidx]
			if exists {
				h.(*HashSet).data[gidx].Or(hval, gval)
			} else {
				hval = big.NewInt(0)
				hval.Set(gval)
				h.(*HashSet).data[gidx] = hval
			}
		}
	default:
	}
}
