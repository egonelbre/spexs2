package set

type Set interface {
	Add(val int)
	Len() int
	Unpack() []int
	Iter(fn func(v int))
}

type Multi interface {
	Set
	AddSet(val Set)
}
