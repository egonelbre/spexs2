package set

import (
	"math/rand"
	"runtime"
	"testing"
)

var goodArr = [...]int{0, 1, 2, 100, 401, 412, 450, 5102, 45104, 51451245}

func mapFrom(arr []int) map[int]int {
	res := make(map[int]int)
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

	for _, val := range set.Iter() {
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
}

func addValues(set Set, n int) {
	for i := range n {
		set.Add(i)
	}
}

func iterate(set Set) {
	sum := 0
	for _, x := range set.Iter() {
		sum += x
	}
}

func testMemoryUse(set Set, n int, t *testing.T) {
	runtime.GC()

	before := new(runtime.MemStats)
	runtime.ReadMemStats(before)

	rng := rand.New(rand.NewSource(5))

	last := 0
	for range n {
		last += 10 + rng.Intn(20)
		set.Add(last)
	}

	after := new(runtime.MemStats)
	runtime.ReadMemStats(after)

	t.Errorf("memory difference %v", after.Alloc-before.Alloc)

	_ = set
}
