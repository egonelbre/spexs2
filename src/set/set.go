package set

type Set interface {
	Add(val int)
	Len() int
	Iter() []int // return unpacked data array
}

type Sorted interface {
	Set
	IsSorted()
}

type Multi interface {
	Set
	AddSet(val Set)
}