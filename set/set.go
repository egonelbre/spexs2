package set

// Set implements an integer set.
type Set interface {
	// Add adds a value to the set.
	Add(val int)
	// Len returns the number of items in the set.
	Len() int
	// All returns all values in the set as a slice.
	All() []int
}

// Multi implements a collection of sets.
type Multi interface {
	Set

	// AddSet adds a new collection to the set.
	AddSet(val Set)
}
