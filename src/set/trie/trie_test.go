package trie

import (
	"bit"
	"testing"
	"testing/quick"
)

func TestDecompose(t *testing.T) {
	check := func(v uint) bool {
		maj, min, bits := decompose(v)
		idx := bit.ScanLeft64(uint64(bits))
		r := compose(maj, min, uint(idx))
		return v == r
	}

	if err := quick.Check(check, nil); err != nil {
		t.Error(err)
	}

	testValues := [...]uint{5102, 412, 51451245, 5102}

	for _, v := range testValues {
		if !check(v) {
			t.Errorf("decompose failed with : %v", v)
		}
	}
}
