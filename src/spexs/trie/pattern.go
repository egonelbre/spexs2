package trie

import (
	"spexs"
	"stats"
)

type Pattern struct {
	Char    Char
	Parent  *Pattern
	Pos     spexs.Set
	IsGroup bool
	IsStar  bool
	count   []int
	occs    []int
	length  int
}

func NewPattern(char Char, parent *Pattern) *Pattern {
	p := &Pattern{}
	p.Char = char
	p.Parent = parent
	if parent != nil {
		p.Pos = spexs.NewHashSet(parent.Pos.Len() / 2)
	} else {
		p.Pos = spexs.NewHashSet(0)
	}
	p.IsGroup = false
	p.IsStar = false
	p.count = make([]int, 0)
	p.occs = make([]int, 0)

	p.length = -1
	return p
}

func NewFullPattern(ref *Reference) *Pattern {
	p := NewPattern(0, nil)
	p.Pos = NewFullSet(ref)
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
	if len(n.count) <= 0 {
		n.count = make([]int, len(ref.Groupings))

		for idx := range n.Pos.Iter() {
			seq := ref.Seqs[idx]
			n.count[seq.Group] += seq.Count
		}
	}
	return n.count[group]
}

func (n *Pattern) Occs(ref *Reference, group int) int {
	if len(n.occs) <= 0 {
		n.occs = make([]int, len(ref.Groupings))

		for idx, mpos := range n.Pos.Iter() {
			seq := ref.Seqs[idx]
			ocs := stats.BitCountInt(mpos)
			n.occs[seq.Group] += seq.Count * ocs
		}
	}
	return n.occs[group]
}
