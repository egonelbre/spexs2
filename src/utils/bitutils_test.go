package utils

import (
	"math/big"
	"testing"
)

var bitCountTests = []struct {
	in  uint
	out int
}{
	{0, 0},
	{1, 1},
	{2, 1},
	{3, 2},
	{4, 1},
	{0xabc, 7},
	{0x80, 1},
	{0x800, 1},
	{0x811, 3},
	{0xfff, 12},
}

func TestBitCount(t *testing.T) {
	for i, test := range bitCountTests {
		if n := BitCount(test.in); n != test.out {
			t.Errorf("#%d got %d want %d", i, n, test.out)
		}
	}
}

func TestBitCount64(t *testing.T) {
	for i, test := range bitCountTests {
		if n := BitCount64(uint64(test.in)); n != test.out {
			t.Errorf("#%d got %d want %d", i, n, test.out)
		}
	}
}

var largeBitCountTests = []struct {
	in  int64
	out int
}{
	{0, 0},
	{1, 1},
	{2, 1},
	{3, 2},
	{4, 1},
	{0xabc, 7},
	{0x80, 1},
	{0x800, 1},
	{0x811, 3},
	{0xfff, 12},
	{0xfff << 32, 12},
	{0xfff<<50 + 3, 14},
}

func TestBitCountInt(t *testing.T) {
	for i, test := range largeBitCountTests {
		v := big.NewInt(test.in)
		if n := BitCountInt(v); n != test.out {
			t.Errorf("#%d got %d want %d", i, n, test.out)
		}
	}
}