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

type Sequence struct {
	Pat   []byte
	Len   int // this refers to rune count in Pat
	Group int // validation or reference group
	Count int
}

type Reference struct {
	Seqs      []Sequence
	Alphabet  []Char
	Groups    map[Char]Group
	Groupings []int
}

func NewReference(size int) *Reference {
	ref := &Reference{}
	ref.Seqs = make([]Sequence, 0, size)
	ref.Alphabet = make([]Char, 0, 8)
	ref.Groups = make(map[Char]Group)
	ref.Groupings = make([]int, 2) // fix use multiple
	return ref
}

func (ref *Reference) Next(idx int, pos byte) (Char, byte, bool) {
	if int(pos) >= len(ref.Seqs[idx].Pat) {
		return 0, 0, false
	}

	rune, width := utf8.DecodeRune(ref.Seqs[idx].Pat[pos:])
	return Char(rune), byte(pos + byte(width)), true
}

func (ref *Reference) ReplaceGroups(pat string) string {
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

func (ref *Reference) AddGroup(group Group) {
	ref.Groups[group.Id] = group
}

func (ref *Reference) AddSequence(pat Sequence) {
	ref.Seqs = append(ref.Seqs, pat)
	ref.Groupings[pat.Group] += 1
}
