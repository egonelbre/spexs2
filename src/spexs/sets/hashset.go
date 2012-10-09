package sets

type HashSet struct {
	data map[int]bool
}

func NewHashSet(size int) *HashSet {
	return &HashSet{make(map[int]bool, size)}
}

func (hs *HashSet) Add(val int) {
	hs.data[val] = true
}

func (hs *HashSet) Contains(val int) bool {
	exists, ok := hs.data[val]
	return exists && ok
}

func (hs *HashSet) Len() int {
	return len(hs.data)
}

func (hs *HashSet) Iter() chan int {
	ch := make(chan int, 100)
	go func(hs *HashSet, ch chan int) {
		for val, ex := range hs.data {
			if ex {
				ch <- val
			}
		}
		close(ch)
	}(hs, ch)
	return ch
}

func (hs *HashSet) Clear() {
	hs.data = nil
}

func (hs *HashSet) AddSet(g *HashSet) {
	for gidx, _ := range g.data {
		hs.data[gidx] = true
	}
}
