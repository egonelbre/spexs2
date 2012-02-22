package spexs

type Char rune

type Group struct {
	Id    Char
	Chars []Char
}

func NewGroup(id Char, chars []Char) *Group {
	return &Group{id, chars}
}

type TrieNode struct {
	Char   Char
	Parent *TrieNode
	Pos    Set
	IsStar bool
}

func NewTrieNode(char Char, parent *TrieNode) *TrieNode {
	return &TrieNode{char, parent, NewHashSet(), false}
}

func NewRootNode(ref Reference) *TrieNode {
	n := &TrieNode{0, nil, NewFullSet(ref), false}
	return n
}

func (n TrieNode) String() string {
	if n.Parent != nil {
		if n.IsStar {
			return n.Parent.String() + string('*') + string(n.Char)
		} else {
			return n.Parent.String() + string(n.Char)
		}
	}
	return "";
}

func TrieCountFilter(limit int) PatternFilter {
	return func(p Pattern) bool {
		return p.(TrieNode).Pos.Length() > limit
	}
}
