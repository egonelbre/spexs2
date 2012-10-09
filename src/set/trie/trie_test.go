package trie

import (
	"bit"
	"testing"
	"testing/quick"
)

func TestDecompose(t *testing.T) {
	check := func(v int) bool {
		maj, min, bits := decompose(v)
		idx := bit.ScanLeft(uint(bits))
		r := compose(maj, min, idx)
		return v == r
	}

	if err := quick.Check(check, nil); err != nil {
		t.Error(err)
	}

	testValues := [...]int{5102, 412, 51451245, 5102}

	for _, v := range testValues {
		if !check(v) {
			t.Errorf("decompose failed with : %v", v)
		}
	}
}
