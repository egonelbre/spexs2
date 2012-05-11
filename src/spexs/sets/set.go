package sets

type BitVector uint16
type Positions map[int]BitVector

type Set interface {
	Add(idx int, pos uint)
	Contains(idx int, pos uint) bool
	Len() int
	Iter() Positions
	Clear()
}
