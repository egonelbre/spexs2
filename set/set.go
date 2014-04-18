package set

type Set interface {
	Add(val int)
	Len() int
	Iter() []int // return unpacked data array
}

type Multi interface {
	Set
	AddSet(val Set)
}
