package spexs

import "utf8"

type ReferencePattern struct {
	Pat []byte
	Count int // this refers to rune count in Pat
}

type UnicodeReference struct {
  Pats []ReferencePattern
  Star Char
  Alphabet []Char
  Groups []Group
}

func (ref UnicodeReference) Next(p Pos) (Char, Pos, bool) {
  idx, pos := PosDecode(p)

  if pos >= len(ref.Pats[idx].Pat) {
  	return 0, EmptyPos, false
  }

  rune, width := utf8.DecodeRune( ref.Pats[idx][pos:] )
  return rune, p + width, true
}
