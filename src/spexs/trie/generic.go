package trie

// this specialises spexs.go implementation

type Pattern *struct {
	*TrieNode
}

type Reference *struct {
	*UnicodeReference
}

type Char rune