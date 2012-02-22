package spexs

import "testing"

func TestEncodeDecode(t *testing.T) {
	
	test := func(idx int, pos byte){
		p := PosEncode(idx, pos)
		didx, dpos := PosDecode(p)
		if idx != didx || pos != dpos {
			t.Errorf("decoding (%v,%v), got (%v,%v)", idx, pos, didx, dpos)
		}
	}

	test(0,0)
	test(3,10)
	test(10,10)
	test(212355,3)
	test(2123,6)
}