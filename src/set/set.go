package set

type Set interface {
	Add(val int)
	Contains(val int) bool
	Len() int
	Iter() []int // return unpacked data array

	// Pack()        // pack internal data
	// Unpack()		 // unpack internal data	
}
