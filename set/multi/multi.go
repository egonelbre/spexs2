package multi

import "github.com/egonelbre/spexs2/set"

type Set struct {
	sets []set.Set
	data []int
}

func New() *Set {
	return &Set{make([]set.Set, 0), nil}
}

func (multi *Set) Add(val int) {
	panic("can't add a single value to multiset")
}

func (multi *Set) AddSet(val set.Set) {
	multi.sets = append(multi.sets, val)
}

func (multi *Set) Len() int {
	c := 0
	for _, s := range multi.sets {
		c += s.Len()
	}
	return c
}

func (multi *Set) unpack() {
	if multi.data != nil {
		return
	}
	sets := make([][]int, 0, len(multi.sets))
	for _, s := range multi.sets {
		sets = append(sets, s.Unpack())
	}
	multi.data = set.MergeSortedInts(sets...)
}

func (multi *Set) Unpack() []int {
	multi.unpack()
	return multi.data
}

func (multi *Set) Iter(fn func(val int)) {
	multi.unpack()
	for _, v := range multi.data {
		fn(v)
	}
}
