package spexs

import "testing"

type p struct {
	idx int
	pos int
}

func TestHashSet(t *testing.T) {
	hs := NewHashSet(10)

	good := map[p]bool{
		p{0, 0}:       true,
		p{13, 1}:      true,
		p{14, 2}:      true,
		p{1025235, 3}: true,
		p{2000, 15}:   true,
		p{200000, 30}: true,
		p{200000, 100}: true,
		p{200000, 3000}: true,
	}

	for pos := range good {
		hs.Add(pos.idx, pos.pos)
	}

	for pos := range good {
		if !hs.Contains(pos.idx, pos.pos) {
			t.Errorf("HashSet didn't contain %v,%v", pos.idx, pos.pos)
		}
	}

	invalid := map[p]bool{
		p{10, 0}:       true,
		p{15, 0}:       true,
		p{20, 3}:       true,
		p{102535, 53}:  true,
		p{2001, 63}:    true,
		p{2000300, 41}: true,
	}

	for pos := range invalid {
		if hs.Contains(pos.idx, pos.pos) {
			t.Errorf("HashSet contained %v,%v", pos.idx, pos.pos)
		}
	}
}

func InsertAndIterate(s Set, patterns int, positions int) {
	for idx := 0; idx < patterns; idx += 1 {
		for pos := 0; pos < positions; pos += 1 {
			s.Add(idx, pos)
		}
	}

	sum := 0
	for x := range s.Iter() {
		sum += x
	}
}

func BenchmarkHashSet(b *testing.B) {
	hs := NewHashSet(10)
	InsertAndIterate(hs, 1000000, 10)
}
