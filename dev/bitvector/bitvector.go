package main

import "fmt"

// use unsafe instead

type Elem struct {
	H uint16 // High bits of addr
	L uint16 // Low bits of addr
	B uint16 // Bit vector bytes
}

type BitVector struct {
	Elems []Elem
}

func New(len int) *BitVector {
	return &BitVector{make([]Elem, len)}
}

func ifthen(cond bool, a int, b int) int {
	if cond {
		return a
	}
	return b
}

func (a *BitVector) And(b *BitVector) *BitVector {
	ai, al := 0, len(a.Elems)
	bi, bl := 0, len(b.Elems)

	rl := ifthen(al < bl, al, bl)
	r, ri := New(rl), 0

	for ai < al && bi < al {
		ae, be := a.Elems[ai], b.Elems[bi]

		switch {
		case ae.H == be.H && ae.L == be.L:
			r.Elems[ri] = Elem{ae.H, ae.L, ae.B & be.B}
			ai += 1
			bi += 1
			ri += 1
		case ae.H <= be.H && ae.L <= be.L:
			ai += 1
		case ae.H >= be.H && ae.L >= be.L:
			bi += 1
		}
	}

	r.Elems = r.Elems[0:ri]
	return r
}

func (b *BitVector) Append(addr uint, bytes uint16) {
	b.Elems = append(b.Elems, Elem{uint16(addr >> 16), uint16(addr & 0xFFFF), bytes})
}

func main() {
	a := New(0)
	a.Append(0, 5)
	a.Append(1, 7)
	a.Append(5, 9)
	a.Append(7, 3)

	b := New(0)
	b.Append(0, 3)
	b.Append(1, 7)
	b.Append(6, 8)
	b.Append(7, 7)

	r := a.And(b)
	fmt.Printf("%+v\n", r)
}
