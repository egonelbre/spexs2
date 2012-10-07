package main

import "fmt"
import "runtime"

// maximal value is 2^(32 + 16 + 4)
// max 

type majkey uint32
type minkey uint16
type bitmap uint16

type Set struct {
	root map[majkey]map[minkey]bitmap
}

func decompose(value uint64) (maj majkey, min minkey, bits bitmap) {
	bits = bitmap(1 << (value & 0xF))   // 4 bits
	min = minkey((value >> 4) & 0xFFFF) // 16 bits
	maj = majkey(value >> 20)           // +4 bytes
	return
}

func decompose1(value uint64) (maj majkey, min minkey, bits bitmap) {
	min = minkey(value & 0xFFFF)              // 16 bits
	bits = bitmap(1 << ((value >> 16) & 0xF)) // 4 bits
	maj = majkey(value >> 20)                 // +4 bytes
	return
}

func New() *Set {
	return &Set{make(map[majkey]map[minkey]bitmap)}
}

func (s *Set) Add(value uint64) {
	maj, min, bits := decompose(value)

	first, exists := s.root[maj]
	if !exists {
		first = make(map[minkey]bitmap)
		s.root[maj] = first
	}

	first[min] |= bits
}

func (s *Set) Union(t *Set) {
	for tmaj, tchild := range t.root {
		schild, ok := s.root[tmaj]
		if !ok {
			schild := make(map[minkey]bitmap)
			s.root[maj] = schild
		}
	}
}

func (s *Set) Print() {
	fmt.Printf("[ %v ]\n", len(s.root))
	for maj, first := range s.root {
		k := 0
		fmt.Printf("  %v\t> %v\n", maj, len(first))
		for min, bits := range first {
			fmt.Printf("    \t %v \t| %b\n", min, bits)
			k += 1
			if k > 5 {
				fmt.Printf("    \t[trunc]\n")
				break
			}
		}
	}
}

func main() {
	//*
	s := New()
	for i := 0; i < 0xFFFFF; i += 1 {
		s.Add(uint64(i))
	}
	stats := new(runtime.MemStats)
	runtime.ReadMemStats(stats)
	fmt.Printf("Alloc %v\n", stats.Alloc)
	s.Print()
}
