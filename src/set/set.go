package set

type Set interface {
	Add(val uint)
	Contains(val uint) bool
	Len() int
	Iter() chan uint
}
