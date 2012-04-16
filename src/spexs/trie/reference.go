package trie

import (
	"bytes"
	"unicode/utf8"
)

type Group struct {
	Id    Char
	Long  string
	Chars []Char
}

func NewGroup(name string, id Char, chars []Char) *Group {
	return &Group{id, name, chars}
}

type ReferencePattern struct {
	Pat   []byte
	Count int // this refers to rune count in Pat
	Group int // validation or reference group
}

type UnicodeReference struct {
	Pats      []ReferencePattern
	Alphabet  []Char
	Groups    map[Char]Group
	Groupings []int
}

func NewUnicodeReference(size int) *UnicodeReference {
	ref := &UnicodeReference{}
	ref.Pats = make([]ReferencePattern, 0, size)
	ref.Alphabet = make([]Char, 0, 8)
	ref.Groups = make(map[Char]Group)
	ref.Groupings = make([]int, 2) // fix use multiple
	return ref
}

func (ref *UnicodeReference) Next(idx int, pos byte) (Char, byte, bool) {
	if int(pos) >= len(ref.Pats[idx].Pat) {
		return 0, 0, false
	}

	rune, width := utf8.DecodeRune(ref.Pats[idx].Pat[pos:])
	return Char(rune), byte(pos + byte(width)), true
}

func (ref *UnicodeReference) ReplaceGroups(pat string) string {
	buf := bytes.NewBufferString("")
	for _, c := range pat {
		if grp, exists := ref.Groups[Char(c)]; exists {
			buf.WriteString(grp.Long)
			continue
		}
		buf.WriteRune(c)
	}
	return string(buf.Bytes())
}

func (ref *UnicodeReference) AddGroup(group Group) {
	ref.Groups[group.Id] = group
}

func (ref *UnicodeReference) AddPattern(pat ReferencePattern) {
	ref.Pats = append(ref.Pats, pat)
	ref.Groupings[pat.Group] += 1
}
