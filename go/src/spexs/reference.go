package spexs

import (
	"unicode/utf8"
	"bytes"
)


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
	Groups   map[Char]Group
}

func NewUnicodeReference(size int) *UnicodeReference{
	ref := &UnicodeReference{}
	ref.Pats = make([]ReferencePattern, 0, size)
	ref.Alphabet = make([]Char, 0, 8)
	ref.Groups = make(map[Char]Group)
	return ref
}

func (ref *UnicodeReference) Next(idx int, pos byte) (Char, byte, bool) {
	if int(pos) >= len(ref.Pats[idx].Pat) {
		return 0, 0, false
	}

	rune, width := utf8.DecodeRune(ref.Pats[idx].Pat[pos:])
	return Char(rune), byte(pos + byte(width)), true
}

func (ref *UnicodeReference) ReplaceGroups(pat string) (string) {
	buf := bytes.NewBufferString("")
	for _, c := range(pat) {
		if grp, exists := ref.Groups[Char(c)]; exists {
			buf.WriteString(grp.Name)
			continue
		}
		buf.WriteRune(c)
	}
	return string(buf.Bytes())
}