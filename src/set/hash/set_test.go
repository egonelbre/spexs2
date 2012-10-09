package hash

import "testing"

func TestRoutine(t *testing.T) {
	hs := New(10)
	good := [...]int{0, 1, 2, 100, 401, 412, 450, 5102, 45104, 51451245}

	for _, val := range good {
		hs.Add(val)
	}

	for _, val := range good {
		if !hs.Contains(val) {
			t.Errorf("HashSet didn't contain %v", val)
		}
	}

	if hs.Len() != len(good) {
		t.Errorf("HashSet had different size %+v, %+v", hs.data, good)
	}

	invalid := [...]int{3, 4, 5, 400, 402, 413, 451, 449, 5101, 5103, 45103}
	for _, val := range invalid {
		if hs.Contains(val) {
			t.Errorf("HashSet contained %v", val)
		}
	}
}

func BenchmarkSet(b *testing.B) {
	hs := New(10)

	for i := 0; i < b.N; i++ {
		hs.Add(i)
	}

	sum := 0
	for x := range hs.Iter() {
		sum += x
	}
}
