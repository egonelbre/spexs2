package hash

type Set struct {
	data map[uint]bool
}

func New(size uint) *Set {
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

func (set *Set) Iter() chan uint {
	ch := make(chan uint, 100)
	go func(set *Set, ch chan uint) {
		for val := range set.data {
			ch <- val
		}
		close(ch)
	}(set, ch)
	return ch
}

func (set *Set) AddSet(other *Set) {
	for val, _ := range other.data {
		set.data[val] = true
	}
}
