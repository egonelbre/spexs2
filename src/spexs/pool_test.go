package spexs

import "testing"

func testTake(t *testing.T, p Pooler, expected string, expectedOk bool) {
	val, ok := p.Take()
	if ok != expectedOk {
		t.Errorf("didn't get correct ok value, got='%v', expected='%v', str='%s'", ok, expectedOk, expected)
	}
	if ok && val.String() != expected {
		t.Errorf("didn't get correct value, got='%s', expected='%s'", val, expected)
	}
}

func TestFifoPool(t *testing.T) {
	p := NewFifoPool()

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

func TestPriorityPool(t *testing.T) {
	lenFitness := func(p *Pattern) float64 {
		return float64(len(p.String()))
	}
	p := NewPriorityPool(lenFitness, 100, true)

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
