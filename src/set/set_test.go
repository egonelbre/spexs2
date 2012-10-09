package set

import (
	"math/rand"
	"runtime"
	"testing"
)

var goodArr = [...]uint{0, 1, 2, 100, 401, 412, 450, 5102, 45104, 51451245}

func mapFrom(arr []uint) map[uint]int {
	res := make(map[uint]int)
	for _, val := range arr {
		res[val] = 0
	}
	return res
}

func testUse(set Set, t *testing.T) {
	good := mapFrom(goodArr[:])

	for val := range good {
		set.Add(val)
	}

	for val := range good {
		if !set.Contains(val) {
			t.Errorf("didn't contain %v", val)
		}
	}

	for val := range set.Iter() {
		it, exists := good[val]
		if !exists {
			t.Errorf("contained value %v that was not added", val)
		}
		good[val] = it + 1
	}

	for val, it := range good {
		if it != 1 {
			t.Errorf("iterated value %v %v times", val, it)
		}
	}

	invalid := [...]uint{3, 4, 5, 400, 402, 413, 451, 449, 5101, 5103, 45103}
	for _, val := range invalid {
		if set.Contains(val) {
			t.Errorf("contained %v", val)
		}
	}
}

func addValues(set Set, n int) {
	for i := uint(0); i < uint(n); i++ {
		set.Add(i)
	}
}

func iterate(set Set) {
	sum := uint(0)
	for x := range set.Iter() {
		sum += x
	}
}

func testMemoryUse(set Set, n int, t *testing.T) {
	runtime.GC()

	before := new(runtime.MemStats)
	runtime.ReadMemStats(before)

	rand.Seed(5)
	last := uint(0)
	for i := 0; i < n; i += 1 {
		last += 10 + uint(rand.Intn(20))
		set.Add(last)
	}

	after := new(runtime.MemStats)
	runtime.ReadMemStats(after)

	t.Errorf("memory difference %v", after.Alloc-before.Alloc)

	_ = set
}
