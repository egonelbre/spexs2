package pool

import (
	. "spexs"
	"testing"
	"unicode/utf8"
)

func pat(s string) *Query {
	r := NewQuery(nil, RegToken{})
	return add(r, s)
}

func add(base *Query, s string) *Query {
	if len(s) <= 0 {
		return base
	}
	rune, size := utf8.DecodeRuneInString(s)
	token := RegToken{Token(rune), false, rune == 'X'}
	n := NewQuery(base, token)
	return add(n, s[size:])
}

func testTake(t *testing.T, p Pooler, expected string, expectedOk bool) {
	val, ok := p.Take()
	if ok != expectedOk {
		t.Errorf("didn't get correct ok value, got='%v', expected='%v', str='%s'", ok, expectedOk, expected)
	}
	if ok && val.StringRaw() != expected {
		t.Errorf("didn't get correct value, got='%s', expected='%s'", val, expected)
	}
}

func TestFifo(t *testing.T) {
	p := NewFifo()

	p.Put(pat("alpha"))
	p.Put(pat("beta"))

	testTake(t, p, "alpha", true)
	testTake(t, p, "beta", true)

	p.Put(pat("gamma"))
	p.Put(pat("delta"))

	testTake(t, p, "gamma", true)
	testTake(t, p, "delta", true)
	testTake(t, p, "", false)
	testTake(t, p, "", false)
}

func TestPriority(t *testing.T) {
	lenFeature := func(q *Query) (float64, string) {
		return float64(q.Len()), ""
	}

	p := NewPriority(lenFeature, 0, true)

	p.Put(pat("bc"))
	p.Put(pat("defg"))
	p.Put(pat("a"))
	p.Put(pat("def"))

	testTake(t, p, "a", true)
	testTake(t, p, "bc", true)

	p.Put(pat("x"))
	p.Put(pat("defgh"))

	testTake(t, p, "x", true)
	testTake(t, p, "def", true)
	testTake(t, p, "defg", true)
	testTake(t, p, "defgh", true)

	testTake(t, p, "", false)
}
