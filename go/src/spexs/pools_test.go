package spexs

import "testing"

type string_pattern struct {
	v string
}

func newString(v string) *string_pattern {
	return &string_pattern{v}
}

func (p *string_pattern) String() string {
	return p.v
}


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
	
	p.Put(newString("alpha"))
	p.Put(newString("beta"))
	
	testTake(t, p, "alpha", true)
	testTake(t, p, "beta", true)
	
	p.Put(newString("gamma"))
	p.Put(newString("delta"))

	testTake(t, p, "gamma", true)
	testTake(t, p, "delta", true)
	testTake(t, p, "", false)
	testTake(t, p, "", false)
}

func TestPriorityPool(t *testing.T) {
	lenFitness := func(p Pattern) float32 {
		return float32(len(p.String()))
	}
	p := NewPriorityPool(lenFitness, 100)
	
	p.Put(newString("bc"))	
	p.Put(newString("defg"))
	p.Put(newString("a"))
	p.Put(newString("def"))
	
	testTake(t, p, "a", true)
	testTake(t, p, "bc", true)

	p.Put(newString("x"))
	p.Put(newString("defgh"))

	testTake(t, p, "x", true)
	testTake(t, p, "def", true)
	testTake(t, p, "defg", true)
	testTake(t, p, "defgh", true)
	
	testTake(t, p, "", false)
}