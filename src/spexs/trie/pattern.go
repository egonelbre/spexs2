package trie

import (
	"spexs"
	"stats"
)

type Pattern struct {
	Char    Char
	Parent  *Pattern
	Pos     spexs.HashSet
	IsGroup bool
	IsStar  bool
	count   [2]int
	occs    [2]int
	length  int
}

func NewPattern(char Char, parent *Pattern) *Pattern {
	p := &Pattern{}
	p.Char = char
	p.Parent = parent
	if parent != nil {
		p.Pos = *spexs.NewHashSet(parent.Pos.Len() / 2)
	} else {
		p.Pos = *spexs.NewHashSet(0)
	}
	p.IsGroup = false
	p.IsStar = false
	p.count[0] = -1
	p.count[1] = -1
	p.occs[0] = -1
	p.occs[1] = -1
	p.length = -1
	return p
}

func NewFullPattern(ref *Reference) *Pattern {
	p := NewPattern(0, nil)
	for idx, pat := range ref.Seqs {
		for i, _ := range pat.Pat {
			p.Pos.Add(idx, i)
		}
	}
	return p
}

func (n *Pattern) String() string {
	if n.Parent != nil {
		if n.IsStar {
			return n.Parent.String() + string('*') + string(n.Char)
		} else {
			return n.Parent.String() + string(n.Char)
		}
	}
	return ""
}

func (n *Pattern) Len() int {
	if n.Parent != nil {
		if n.length < 0 {
			n.length = n.Parent.Len() + 1
		}
		return n.length
	}
	return 0
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
			ocs := stats.BitCountInt(mpos)
			n.occs[seq.Group] += seq.Count * ocs
		}
	}
	return n.occs[group]
}
