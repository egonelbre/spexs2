package spexs

type TrieNode struct {
	Char   Char
	Parent *TrieNode
	Pos    Set
	IsGroup bool
	IsStar bool
	length int
	complexity int
}

func NewTrieNode(char Char, parent *TrieNode) *TrieNode {
	return &TrieNode{char, parent, NewHashSet(), false, false, -1, -1}
}

func NewFullNodeFromRef(ref *UnicodeReference) *TrieNode {
	n := &TrieNode{0, nil, NewFullSet(ref), false, false, -1, -1}
	return n
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

func (n *TrieNode) Length() int {
	if n.Parent != nil {
		if n.length < 0 {
			n.length = n.Parent.Length() + 1
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
