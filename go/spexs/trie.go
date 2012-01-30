package trie

type Char int

type Group struct {
  Id Char
  Chars []Char
}

func NewGroup( id Char, chars []Char ) *Group {
  return &Group{ id, chars }
}

type TrieNode struct {
  Char Char
  Parent *TrieNode
  Pos *Set
  IsStar bool
}

func NewTrieNode(char Char, parent TrieNode) *TrieNode{
  return &TrieNode{char, parent, NewHashSet(), false}
}

func NewRootNode(parent TrieNode, ref Reference) *TrieNode {
	n := &TrieNode{0, nil, NewFullSet(ref), false}
	return n
}

func (n * TrieNode) ToString() string {
	if n.Parent == nil {
		return ""
	} else {
    if n.IsStar {
      return n.Parent.ToString() + '*' +  string( n.Char )
    } else {
      return n.Parent.ToString() + string( n.Char )
    }
		
	}
}

func TrieCountFilter( limit int ) PatternFilter {
  return func(p Pattern) bool {
    return p.(TrieNode).Pos.Length() > limit
  }
}
