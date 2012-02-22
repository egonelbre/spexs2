package spexs

import "unicode/utf8"

type Char rune

type Group struct {
	Id    Char
	Chars []Char
}

func NewGroup(id Char, chars []Char) *Group {
	return &Group{id, chars}
}

type ReferencePattern struct {
	Pat   []byte
	Count int // this refers to rune count in Pat
}

type UnicodeReference struct {
	Pats     []ReferencePattern
	Alphabet []Char
	Groups   []Group
}

func (ref UnicodeReference) Next(p Pos) (Char, Pos, bool) {
	idx, pos := PosDecode(p)

	if int(pos) >= len(ref.Pats[idx].Pat) {
		return 0, EmptyPos, false
	}

	rune, width := utf8.DecodeRune(ref.Pats[idx].Pat[pos:])
	next_pos := Pos(int64(p) + int64(width))
	return Char(rune), next_pos, true
}
