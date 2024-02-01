package pool

import (
	"testing"
	"unicode/utf8"

	"github.com/egonelbre/spexs2/search"
)

func pat(s string) *search.Query {
	r := search.NewQuery(nil, search.RegToken{})
	r.Db = search.NewDatabase()
	return add(r, s)
}

func add(base *search.Query, s string) *search.Query {
	if len(s) == 0 {
		return base
	}
	rune, size := utf8.DecodeRuneInString(s)

	token := search.RegToken{Token: search.Token(rune), Flags: search.IsSingle}
	if rune == 'X' {
		token.Flags = search.IsStar
	}
	n := search.NewQuery(base, token)
	return add(n, s[size:])
}

func testPop(t *testing.T, p search.Pooler, expected string, expectedOk bool) {
	val, ok := p.Pop()
	if ok != expectedOk {
		t.Errorf("didn't get correct ok value, got='%v', expected='%v', str='%s'", ok, expectedOk, expected)
	}
	if ok && val.StringRaw() != expected {
		t.Errorf("didn't get correct value, got='%s', expected='%s'", val, expected)
	}
}

func TestFifo(t *testing.T) {
	p := NewQueue()

	p.Push(pat("alpha"))
	p.Push(pat("beta"))

	testPop(t, p, "alpha", true)
	testPop(t, p, "beta", true)

	p.Push(pat("gamma"))
	p.Push(pat("delta"))

	testPop(t, p, "gamma", true)
	testPop(t, p, "delta", true)
	testPop(t, p, "", false)
	testPop(t, p, "", false)
}

func TestPriority(t *testing.T) {
	lenFeature := func(q *search.Query) (float64, string) {
		return float64(q.Len()), ""
	}

	p := NewPriority([]search.Feature{lenFeature}, 0)

	p.Push(pat("bc"))
	p.Push(pat("defg"))
	p.Push(pat("a"))
	p.Push(pat("def"))

	testPop(t, p, "a", true)
	testPop(t, p, "bc", true)

	p.Push(pat("x"))
	p.Push(pat("defgh"))

	testPop(t, p, "x", true)
	testPop(t, p, "def", true)
	testPop(t, p, "defg", true)
	testPop(t, p, "defgh", true)

	testPop(t, p, "", false)
}
