package set

type Set interface {
	Add(val int)
	Len() int
	Iter() []int // return unpacked data array
	IsSorted() bool
}

type Multi interface {
	Set
	AddSet(val Set)
}
