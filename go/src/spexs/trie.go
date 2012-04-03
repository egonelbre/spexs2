package spexs

type TrieNode struct {
	Char   Char
	Parent *TrieNode
	Pos    Set
	IsGroup bool
	IsStar bool
	length int
	complexity int
	pvalue float64
}

func NewTrieNode(char Char, parent *TrieNode) *TrieNode {
	return &TrieNode{char, parent, NewHashSet(parent.Pos.Len() / 2), false, false, -1, -1, -1}
}

func NewFullNodeFromRef(ref *UnicodeReference) *TrieNode {
	return &TrieNode{0, nil, NewFullSet(ref), false, false, -1, -1, -1}
}

func (n *TrieNode) String() string {
	if n.Parent != nil {
		if n.IsStar {
			return n.Parent.String() + string('*') + string(n.Char)
		} else {
			return n.Parent.String() + string(n.Char)
		}
	}
	return "";
}

func (n *TrieNode) Len() int {
	if n.Parent != nil {
		if n.length < 0 {
			n.length = n.Parent.Len() + 1
		}
		return n.length
	}
	return 0
}

func (n *TrieNode) Complexity() int {
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

func (n *TrieNode) PValue(ref *UnicodeReference) float64 {
	if n.pvalue >= 0 {
		return n.pvalue
	}

	counts := make([]int, len(ref.Groupings))

	for idx := range n.Pos.Iter() {
		counts[ref.Pats[idx].Group] += 1
	}

	n.pvalue = HypergeometricSplitLog(counts[0], counts[1], ref.Groupings[0], ref.Groupings[1])
	return n.pvalue
}
