package trie

import "bit"

type majkey uint32
type minkey uint16
type bitmap uint16

type Set struct {
	root map[majkey]map[minkey]bitmap
}

func decompose(value int) (maj majkey, min minkey, bits bitmap) {
	bits = bitmap(1 << uint(value&0xF)) // 4 bits
	min = minkey((value >> 4) & 0xFFFF) // 16 bits
	maj = majkey(value >> 20)           // +4 bytes
	return
}

func compose(maj majkey, min minkey, idx int) int {
	return (int(maj) << 20) | (int(min) << 4) | idx
}

func New() *Set {
	return &Set{make(map[majkey]map[minkey]bitmap)}
}

func (s *Set) Add(value int) {
	maj, min, bits := decompose(value)

	first, exists := s.root[maj]
	if !exists {
		first = make(map[minkey]bitmap)
		s.root[maj] = first
	}

	first[min] |= bits
}

func (s *Set) Contains(value int) bool {
	maj, min, bits := decompose(value)
	mmin, exists := s.root[maj]
	if !exists {
		return false
	}
	return mmin[min]&bits != 0
}

func (s *Set) Iter() chan int {
	ch := make(chan int, 100)
	go func(s *Set, ch chan int) {
		for maj, mmin := range s.root {
			for min, bits := range mmin {
				for k := 0; k < 16; k += 1 {
					if (bits>>uint(k))&1 == 1 {
						ch <- compose(maj, min, k)
					}
				}
			}
		}

		close(ch)
	}(s, ch)
	return ch
}

func (s *Set) Len() int {
	count := 0
	for _, m := range s.root {
		for _, v := range m {
			count += bit.Count64(uint64(v))
		}
	}
	return count
}

func (s *Set) AddSet(t *Set) {
	/*for tmaj, tchild := range t.root {
		schild, ok := s.root[tmaj]
		if !ok {
			schild := make(map[minkey]bitmap)
			s.root[tmaj] = schild
		}
	}*/
}
