package bin

type Set struct {
	data map[uint]bitvector
}

type bitvector uint64

const bitsSize = 64
const bitsCount = 6
const bitsOffset = 4
const offsetMask = (1 << bitsOffset) - 1
const bitsMask = (1 << bitsCount) - 1

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

func New(size int) *Set {
	return &Set{make(map[uint]bitvector, size)}
}

func (set *Set) Add(val uint) {
	idx, bits := decompose(val)
	set.data[idx] |= bits
}

func (set *Set) Contains(val uint) bool {
	idx, bits := decompose(val)
	return set.data[idx]&bits != 0
}

func (set *Set) Len() int {
	return len(set.data)
}

func (set *Set) Iter() chan uint {
	ch := make(chan uint, 100)
	go func(set *Set, ch chan uint) {
		for val, bits := range set.data {
			for k := uint(0); k < bitsSize; k += 1 {
				if (bits>>k)&1 == 1 {
					ch <- compose(val, k)
				}
			}
		}
		close(ch)
	}(set, ch)
	return ch
}

func (set *Set) AddSet(other *Set) {
	for val, bits := range other.data {
		set.data[val] |= bits
	}
}
