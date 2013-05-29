package multi

import (
	"set"
	"sort"
)

type Set struct {
	sets []set.Set
}

func New() *Set {
	return &Set{make([]set.Set, 0)}
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

func (multi *Set) Iter() []int {
	sets := make([][]int, 0, len(multi.sets))
	sorted := true
	for _, s := range multi.sets {
		if _, ok := s.(set.SortedSet); !ok {
			sorted = false
		}
		sets = append(sets, s.Iter())
	}

	if sorted {
		return set.MergeSortedUniqueInts(sets...)
	}

	result := make([]int, 0)
	for _, data := range sets {
		result = append(result, data...)
	}

	return result
}

var _ = sort.Ints
