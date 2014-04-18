package trie

import "bit"

type majkey uint32
type minkey uint16
type bitmap uint16

type Set struct {
	root map[majkey]map[minkey]bitmap
}

func decompose(value uint) (maj majkey, min minkey, bits bitmap) {
	bits = bitmap(1 << uint(value&0xF)) // 4 bits
	min = minkey((value >> 4) & 0xFFFF) // 16 bits
	maj = majkey(value >> 20)           // +4 bytes
	return
}

func compose(maj majkey, min minkey, idx uint) uint {
	return uint(maj)<<20 | uint(min)<<4 | idx
}

func New() *Set {
	return &Set{make(map[majkey]map[minkey]bitmap)}
}

func (set *Set) Add(value int) {
	maj, min, bits := decompose(uint(value))

	first, exists := set.root[maj]
	if !exists {
		first = make(map[minkey]bitmap)
		set.root[maj] = first
	}

	first[min] |= bits
}

func (set *Set) Contains(value int) bool {
	maj, min, bits := decompose(uint(value))
	mmin, exists := set.root[maj]
	if !exists {
		return false
	}
	return mmin[min]&bits != 0
}

func (set *Set) Iter() []int {
	iter := make([]int, set.Len())
	i := 0
	for maj, mmin := range set.root {
		for min, bits := range mmin {
			for k := uint(0); bits > 0; k += 1 {
				if bits&(1<<k) > 0 {
					iter[i] = int(compose(maj, min, k))
					bits = bits &^ (1 << k)
					i += 1
				}
			}
		}
	}
	return iter
}

func (set *Set) Len() int {
	count := 0
	for _, m := range set.root {
		for _, v := range m {
			count += bit.Count64(uint64(v))
		}
	}
	return count
}

func (set *Set) AddSet(other *Set) {
	for maj, mmin := range other.root {
		mm, exists := set.root[maj]
		if !exists {
			mm = make(map[minkey]bitmap)
			set.root[maj] = mm
		}
		for min, bits := range mmin {
			mm[min] |= bits
		}
	}
}

func (set *Set) IsSorted() bool {
	return false
}
