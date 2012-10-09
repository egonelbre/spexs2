package hash

type Set struct {
	data map[int]bool
}

func New(size int) *Set {
	return &Set{make(map[int]bool, size)}
}

func (hs *Set) Add(val int) {
	hs.data[val] = true
}

func (hs *Set) Contains(val int) bool {
	exists, ok := hs.data[val]
	return exists && ok
}

func (hs *Set) Len() int {
	return len(hs.data)
}

func (hs *Set) Iter() chan int {
	ch := make(chan int, 100)
	go func(hs *Set, ch chan int) {
		for val, ex := range hs.data {
			if ex {
				ch <- val
			}
		}
		close(ch)
	}(hs, ch)
	return ch
}

func (hs *Set) Clear() {
	hs.data = nil
}

func (hs *Set) AddSet(g *Set) {
	for val, _ := range g.data {
		hs.data[val] = true
	}
}
