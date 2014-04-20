package search

import (
	"testing"
	"unicode/utf8"
)

func pat(s string) *Query {
	r := NewQuery(nil, RegToken{})
	r.Db = NewDatabase()
	return add(r, s)
}

func add(base *Query, s string) *Query {
	if len(s) <= 0 {
		return base
	}
	rune, size := utf8.DecodeRuneInString(s)

	tok := RegToken{Token(rune), IsSingle}
	if rune == 'X' {
		tok.Flags = IsStar
	}
	n := NewQuery(base, tok)
	return add(n, s[size:])
}

func TestTrieNodeString(t *testing.T) {
	root := pat("")

	a := add(root, "hello")
	b := add(a, " world")
	c := add(root, "hello mXgic")
	d := add(root, "testing")
	e := add(d, " heist")

	test := func(n *Query, s string) {
		if n.StringRaw() != s {
			t.Errorf("wrong result got='%s' expected='%s'", n.String(), s)
		}
	}

	test(a, "hello")
	test(b, "hello world")
	test(c, "hello m*Xgic")
	test(d, "testing")
	test(e, "testing heist")
}
