package spexs

import "unicode/utf8"

type Char rune

type Group struct {
	Id    Char
	Name  string
	Chars []Char
}

func NewGroup(name string, id Char, chars []Char) *Group {
	return &Group{id, name, chars}
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

func (ref *UnicodeReference) Next(idx int, pos int) (Char, int, bool) {
	if int(pos) >= len(ref.Pats[idx].Pat) {
		return 0, 0, false
	}

	rune, width := utf8.DecodeRune(ref.Pats[idx].Pat[pos:])
	return Char(rune), pos + width, true
}
