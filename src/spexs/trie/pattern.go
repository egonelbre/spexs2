package trie

import (
	"spexs"
	"stats"
)

type Pattern struct {
	Char       Char
	Parent     *Pattern
	Pos        spexs.Set
	IsGroup    bool
	IsStar     bool
	length     int
	complexity int
	pvalue     float64
	ng int
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
	p.ng = -1
	p.length = -1
	p.complexity = -1
	p.pvalue = -1
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

func (n *Pattern) NG() int {
	if n.Parent != nil {
		if n.ng < 0 {
			if n.IsGroup || n.IsStar {
				n.ng = n.Parent.NG()
			} else {
				n.ng = 1 + n.Parent.NG()
			}
			
		}
		return n.ng
	}
	return 0
}


func (n *Pattern) Complexity() int {
	if n.Parent != nil {
		if n.complexity < 0 {
			if n.IsStar {
				n.complexity = n.Parent.Complexity() + 4
			} else if n.IsGroup {
				n.complexity = n.Parent.Complexity() + 2
			} else {
				n.complexity = n.Parent.Complexity() + 1
			}
		}
		return n.complexity
	}
	return 0
}

func (n *Pattern) PValue(ref *Reference) float64 {
	if n.pvalue >= 0 {
		return n.pvalue
	}

	counts := make([]int, len(ref.Groupings))

	for idx := range n.Pos.Iter() {
		counts[ref.Seqs[idx].Group] += 1
	}

	n.pvalue = stats.HypergeometricSplit(
		counts[0], counts[1], 
		ref.Groupings[0], ref.Groupings[1])
	return n.pvalue
}
