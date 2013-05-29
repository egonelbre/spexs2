package set

type Set interface {
	Add(val int)
	Len() int
	Iter() []int // return unpacked data array

	// Pack()        // pack internal data
	// Unpack()		 // unpack internal data
}

type SortedSet interface {
	Set
	IsSorted()
}

type MultiSet interface {
	Set
	AddSet(val Set)
}