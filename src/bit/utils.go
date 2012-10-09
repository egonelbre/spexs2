package bit

import "math/big"

const (
	m1  = 0x5555555555555555 //binary: 0101...
	m2  = 0x3333333333333333 //binary: 00110011..
	m4  = 0x0f0f0f0f0f0f0f0f //binary:  4 zeros,  4 ones ...
	m8  = 0x00ff00ff00ff00ff //binary:  8 zeros,  8 ones ...
	m16 = 0x0000ffff0000ffff //binary: 16 zeros, 16 ones ...
	m32 = 0x00000000ffffffff //binary: 32 zeros, 32 ones

	b32  = 0xFFFFFFFF
	ms1  = m1 & b32
	ms2  = m2 & b32
	ms4  = m4 & b32
	ms8  = m8 & b32
	ms16 = m16 & b32
	ms32 = m32 & b32
)

func Count64(x uint64) int {
	x = (x & m1) + ((x >> 1) & m1)
	x = (x & m2) + ((x >> 2) & m2)
	x = (x & m4) + ((x >> 4) & m4)
	x = (x & m8) + ((x >> 8) & m8)
	x = (x & m16) + ((x >> 16) & m16)
	x = (x & m32) + ((x >> 32) & m32)
	return int(x)
}

func Count(x uint) int {
	x = (x & ms1) + ((x >> 1) & ms1)
	x = (x & ms2) + ((x >> 2) & ms2)
	x = (x & ms4) + ((x >> 4) & ms4)
	x = (x & ms8) + ((x >> 8) & ms8)
	x = (x & ms16) + ((x >> 16) & ms16)
	return int(x)
}

func CountInt(x *big.Int) int {
	total := 0
	for _, w := range x.Bits() {
		total += Count64(uint64(w))
	}
	return total
}

func ScanLeft(x uint) int {
	for k := byte(0); k < 32; k += 1 {
		if x&(1<<k) != 0 {
			return int(k)
		}
	}
	return -1
}

func ScanLeft64(x uint64) int {
	for k := byte(0); k < 64; k += 1 {
		if x&(1<<k) != 0 {
			return int(k)
		}
	}
	return -1
}

func ScanLeftInt(x *big.Int) int {
	for k := 0; k < x.BitLen(); k += 1 {
		if x.Bit(k) != 0 {
			return k
		}
	}
	return -1
}
