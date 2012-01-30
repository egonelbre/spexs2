package spexs

import utf8

type ReferencePattern {
	Pat []byte
	Length int
}

type UnicodeReference struct {
  Pats []ReferencePattern
  Star Char
  Alphabet []Char
  Groups []Group
}

func (ref *UnicodeReferences) Next(p Pos) (Char, Pos, bool) {
  idx, pos := PosDecode(p)

  if pos >= ref.Pats[idx].Length {
  	return 0, EmptyPos, false
  }

  rune, width := utf8.DecodeRune( ref.Pats[idx][pos:] )
  return rune, p + width, true
}

func (ref *UnicodeReferences) SeekFrom(p Pos, amount int) (Char, Pos, bool) {
  idx, pos := PosDecode(p)

  /* TODO: seek properly utf8 style */
  if p + amount >= ref.Pats[idx].Length {
  	return 0, EmptyPos, false
  }

  rune := ref.Pats[idx][ pos + amount ]
  return rune, p + amount, true
}
