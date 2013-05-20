package hash

type Set struct {
	data map[uint]struct{}
}

const initSize = 30

func New() *Set {
	return &Set{make(map[int]struct{}, initSize)}
}

func (set *Set) Add(val int) {
	set.data[val] = struct{}{}
}

func (set *Set) Contains(val int) bool {
	exists, ok := set.data[val]
	return exists && ok
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

func (set *Set) AddSet(other *Set) {
	for val, _ := range other.data {
		set.data[val] = true
	}
}
