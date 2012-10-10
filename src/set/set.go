package set

type Set interface {
	Add(val uint)
	Contains(val uint) bool
	Len() int
	Iter() []uint // return unpacked data array

	// Pack()        // pack internal data
	// Unpack()		 // unpack internal data	
}
