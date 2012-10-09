package bin

import (
	"bit"
	"testing"
	"testing/quick"
)

func TestDecompose(t *testing.T) {
	check := func(v uint) bool {
		maj, bits := decompose(v)
		idx := bit.ScanLeft64(uint64(bits))
		r := compose(maj, uint(idx))
		return v == r
	}

	if err := quick.Check(check, nil); err != nil {
		t.Error(err)
	}

	testValues := [...]uint{0, 1, 2, 512, 4999, 5000, 5101, 5102, 5103, 412, 51451245}

	for _, v := range testValues {
		if !check(v) {
			maj, bits := decompose(v)
			idx := bit.ScanLeft64(uint64(bits))
			r := compose(maj, uint(idx))
			t.Errorf("decompose failed with : %v [%x, %x] => [%v]", v, maj, bits, r)
		}
	}
}
