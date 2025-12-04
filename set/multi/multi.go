package multi

import "github.com/egonelbre/spexs2/set"

// Set represents a multiset of integers.
//
// The implementation assumes that the sets added are non-overlapping.
type Set struct {
	sets []set.Set
}

// New creates a new empty multiset.
func New() *Set {
	return &Set{make([]set.Set, 0)}
}

// Add adds a new element to the multiset.
func (multi *Set) Add(val int) {
	panic("can't add a single value to multiset")
}

// AddSet adds a new set to the multiset.
func (multi *Set) AddSet(val set.Set) {
	multi.sets = append(multi.sets, val)
}

// Len returns the number of elements in the multiset.
func (multi *Set) Len() int {
	c := 0
	for _, s := range multi.sets {
		c += s.Len()
	}
	return c
}

// All returns all elements in the multiset.
func (multi *Set) All() []int {
	sets := make([][]int, 0, len(multi.sets))

	for _, s := range multi.sets {
		sets = append(sets, s.All())
	}

	return mergeSortedInts(sets...)
}
