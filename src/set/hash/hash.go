package hash

type Set struct {
	data map[uint]bool
}

func New(size int) *Set {
	return &Set{make(map[uint]bool, size)}
}

func (set *Set) Add(val uint) {
	set.data[val] = true
}

func (set *Set) Contains(val uint) bool {
	exists, ok := set.data[val]
	return exists && ok
}

func (set *Set) Len() int {
	return len(set.data)
}

func (set *Set) Iter() []uint {
	iter := make([]uint, set.Len())
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
