package spexs

import (
	"spexs/sets"
	"utils"
)

type RegChar struct {
	Char    Char
	IsGroup bool
	IsStar  bool
}

type Pattern struct {
	Pat   []RegChar
	Pos   sets.HashSet
	count [2]int
	occs  [2]int
}

func NewPattern(char Char, parent *Pattern, isGroup bool, isStar bool) *Pattern {
	p := &Pattern{}

	if parent != nil {
		p.Pat = append(parent.Pat, RegChar{char, isGroup, isStar})
		p.Pos = *sets.NewHashSet(parent.Pos.Len() / 8)
	} else {
		p.Pat = nil
		p.Pos = *sets.NewHashSet(0)
	}

	p.count[0] = -1
	p.count[1] = -1
	p.occs[0] = -1
	p.occs[1] = -1
	return p
}

func NewFullPattern(ref *Reference) *Pattern {
	p := NewPattern(0, nil, false, false)
	for idx, pat := range ref.Seqs {
		for i, _ := range pat.Pat {
			p.Pos.Add(idx, uint(i))
		}
	}
	return p
}

func (n *Pattern) String() string {
	r := ""
	for _, e := range n.Pat {
		if e.IsStar {
			r += "*" + string(e.Char)
		} else {
			r += string(e.Char)
		}
	}
	return r
}

func (n *Pattern) Len() int {
	return len(n.Pat)
}

func (n *Pattern) Count(ref *Reference, group int) int {
	if n.count[0] < 0 {
		n.count[0] = 0
		n.count[1] = 0

		for idx := range n.Pos.Iter() {
			seq := ref.Seqs[idx]
			n.count[seq.Group] += seq.Count
		}
	}
	return n.count[group]
}

func (n *Pattern) Occs(ref *Reference, group int) int {
	if n.occs[0] < 0 {
		n.occs[0] = 0
		n.occs[1] = 0

		for idx, mpos := range n.Pos.Iter() {
			seq := ref.Seqs[idx]
			ocs := utils.BitCount64(uint64(mpos))
			n.occs[seq.Group] += seq.Count * ocs
		}
	}
	return n.occs[group]
}
