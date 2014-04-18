package hash

type Set struct {
	data map[int]struct{}
}

const initSize = 30

func New() *Set {
	return &Set{make(map[int]struct{}, initSize)}
}

func (set *Set) Add(val int) {
	set.data[val] = struct{}{}
}

func (set *Set) Contains(val int) bool {
	_, ok := set.data[val]
	return ok
}

func (set *Set) Len() int {
	return len(set.data)
}

func (set *Set) Iter() []int {
	iter := make([]int, set.Len())
	i := 0
	for val := range set.data {
		iter[i] = val
		i += 1
	}
	return iter
}

func (set *Set) IsSorted() bool {
	return false
}