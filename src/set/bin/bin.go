package bin

import "bit"

type Set struct {
	data map[uint]bitvector
}

type bitvector uint64

const (
	initSize   = 10
	bitsSize   = 64
	bitsCount  = 6 // log2 bitsSize
	bitsOffset = 10
	offsetMask = (1 << bitsOffset) - 1
	bitsMask   = (1 << bitsCount) - 1
)

func decompose(val uint) (uint, bitvector) {
	high := (val >> (bitsCount + bitsOffset)) << bitsOffset
	low := val & offsetMask
	idx := high | low
	bits := bitvector(1) << ((val >> bitsOffset) & bitsMask)
	return idx, bits
}

func compose(idx uint, pos uint) uint {
	high := uint((idx &^ offsetMask) << bitsCount)
	low := uint(idx & offsetMask)
	mid := uint(pos << bitsOffset)
	return high | low | mid
}

func New() *Set {
	return &Set{make(map[uint]bitvector, initSize)}
}

func (set *Set) Add(val int) {
	idx, bits := decompose(uint(val))
	set.data[idx] |= bits
}

func (set *Set) Contains(val int) bool {
	idx, bits := decompose(uint(val))
	return set.data[idx]&bits != 0
}

func (set *Set) Len() int {
	count := 0
	for _, bits := range set.data {
		count += bit.Count64(uint64(bits))
	}
	return count
}

func (set *Set) Iter() []int {
	iter := make([]int, set.Len())
	i := 0
	for val, bits := range set.data {
		for k := uint(0); bits > 0; k += 1 {
			if bits&(1<<k) > 0 {
				iter[i] = int(compose(val, k))
				bits = bits &^ (1 << k)
				i += 1
			}
		}
	}
	return iter
}

func (set *Set) AddSet(other *Set) {
	for val, bits := range other.data {
		set.data[val] |= bits
	}
}

func (set *Set) IsSorted() bool {
	return false
}