package sets

type BitVector uint64
type Positions map[int]BitVector

type Set interface {
	Add(idx int, pos int)
	Contains(idx int, pos int) bool
	Len() int
	Iter() Positions
	Clear()
}
