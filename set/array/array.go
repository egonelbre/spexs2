package array

// Set represents a set of integers.
//
// It's guaranteed that the elements are added in increasing order
// and there are no duplicates.
type Set []int

// New creates a new empty set.
func New() *Set {
	arr := Set(make([]int, 0, 8))
	return &arr
}

// From creates a new set from the given values.
func From(values ...int) *Set {
	v := Set(values)
	return &v
}

// Add adds a value to the set.
func (s *Set) Add(v int) {
	*s = append(*s, v)
}

// Len returns the number of elements in the set.
func (s *Set) Len() int {
	return len(*s)
}

// All returns all elements in the set.
func (s *Set) All() []int {
	return []int(*s)
}
