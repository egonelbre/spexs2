package spexs

import (
	"testing"
	"unicode/utf8"
)


func TestTrieNodeString(t *testing.T) {
	root := NewTrieNode(0, nil)
	var add func (base *TrieNode, s string) *TrieNode

	add = func (base *TrieNode, s string) *TrieNode {
		if len(s) <= 0 {
			return base
		}
		rune, size := utf8.DecodeRuneInString(s)
		n := NewTrieNode(Char(rune), base)
		if rune == 'x' {
			n.IsStar = true
		}
		return add(n, s[size:])
	}

	a := add(root, "hello")
	b := add(a, " world")
	c := add(root, "hello mxgic")
	d := add(root, "testing")
	e := add(d, " heist")

	test := func(n *TrieNode, s string) {
		if n.String() != s {
			t.Errorf("wrong result got='%s' expected='%s'", n.String(), s)
		}
	}

	test(a, "hello")
	test(b, "hello world")
	test(c, "hello m*xgic")
	test(d, "testing")
	test(e, "testing heist")
}