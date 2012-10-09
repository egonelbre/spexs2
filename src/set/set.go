package set

type Set interface {
	Add(val int)
	Contains(val int) bool
	Len() int
	Iter() chan int
	Clear()
}
