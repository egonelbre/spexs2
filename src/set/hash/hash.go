package hash

type Set struct {
	data map[int]bool
}

func New(size int) *Set {
	return &Set{make(map[int]bool, size)}
}

func (set *Set) Add(val int) {
	set.data[val] = true
}

func (set *Set) Contains(val int) bool {
	exists, ok := set.data[val]
	return exists && ok
}

func (set *Set) Len() int {
	return len(set.data)
}

func (set *Set) Iter() chan int {
	ch := make(chan int, 100)
	go func(set *Set, ch chan int) {
		for val, ex := range set.data {
			if ex {
				ch <- val
			}
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
