package sets

type BitVector uint16
type Positions map[int]BitVector

type Set interface {
	Add(idx int, pos int)
	Contains(idx int, pos int) bool
	Len() int
	Iter() Positions
	Clear()
}
