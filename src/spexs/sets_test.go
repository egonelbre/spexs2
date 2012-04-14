package spexs

import "testing"

func TestEncodeDecode(t *testing.T) {

	test := func(idx int, pos byte) {
		p := PosEncode(idx, pos)
		didx, dpos := PosDecode(p)
		if idx != didx || pos != dpos {
			t.Errorf("decoding (%v,%v), got (%v,%v)", idx, pos, didx, dpos)
		}
	}

	test(0, 0)
	test(3, 10)
	test(10, 10)
	test(212355, 3)
	test(2123, 6)
}

func TestHashSet(t *testing.T) {
	hs := NewHashSet()

	p := PosEncode
	good := map[Pos]bool{
		p(0, 0):        true,
		p(13, 1):       true,
		p(14, 2):       true,
		p(1025235, 3):  true,
		p(2000, 15):    true,
		p(2000000, 50): true,
	}

	for pos := range good {
		hs.Add(pos)
	}

	for pos := range good {
		if !hs.Contains(pos) {
			idx, posv := PosDecode(pos)
			t.Errorf("HashSet didn't contain %v,%v", idx, posv)
		}
	}

	invalid := map[Pos]bool{
		p(10, 0):       true,
		p(15, 0):       true,
		p(20, 3):       true,
		p(102535, 53):  true,
		p(2001, 63):    true,
		p(2000300, 41): true,
	}

	for pos := range invalid {
		if hs.Contains(pos) {
			idx, posv := PosDecode(pos)
			t.Errorf("HashSet contained %v,%v", idx, posv)
		}
	}

	for pos := range hs.Iter() {
		positive, exists := good[pos]
		if exists && positive {
			good[pos] = false
		} else {
			idx, posv := PosDecode(pos)
			if !exists {
				t.Errorf("HashSet shouldn't contain %v,%v", idx, posv)
			} else if !positive {
				t.Errorf("HashSet iterated %v,%v twice", idx, posv)
			}

		}
	}

	for pos, positive := range good {
		if positive {
			idx, posv := PosDecode(pos)
			t.Errorf("HashSet didn't iterate %v,%v", idx, posv)
		}
	}
}

func InsertAndIterate(s Set, patterns int, positions int) {
	for idx := 0; idx < patterns; idx += 1 {
		for pos := 0; pos < positions; pos += 1 {
			p := PosEncode(idx, byte(pos))
			s.Add(p)
		}
	}

	sum := 0
	for x := range s.Iter() {
		_, pos := PosDecode(x)
		sum += int(pos)
	}
}

func BenchmarkHashSet(b *testing.B) {
	hs := NewHashSet()
	InsertAndIterate(hs, 1000000, 10)
}
