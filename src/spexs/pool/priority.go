package pool

import (
	"container/heap"
	. "spexs"
)

type Priority struct {
	token     chan int
	items     []*Pattern
	Fitness   FitnessFunc
	length    int
	limit     int
	ascending bool
}

func NewPriority(fitness FitnessFunc, limit int, ascending bool) *Priority {
	p := &Priority{}
	p.token = make(chan int, 1)
	p.items = make([]*Pattern, limit+100)
	p.length = 0
	p.limit = limit
	p.Fitness = fitness
	p.token <- 1
	p.ascending = ascending

	heap.Init(p)
	return p
}

func (p *Priority) IsEmpty() bool {
	return p.length == 0
}

func (p *Priority) Take() (*Pattern, bool) {
	<-p.token
	if p.IsEmpty() {
		p.token <- 1
		return nil, false
	}
	v := heap.Pop(p)
	p.token <- 1
	return v.(*Pattern), true
}

func (p *Priority) Put(pat *Pattern) {
	<-p.token
	heap.Push(p, pat)
	//if p.limit > 0 && p.Len() > p.limit {
	//	heap.Pop(p)
	//}
	p.token <- 1
}

// sort.Interface
func (p *Priority) Len() int {
	return p.length
}

func (p *Priority) Swap(i, j int) {
	p.items[i], p.items[j] = p.items[j], p.items[i]
}

func (p *Priority) Less(i, j int) bool {
	if p.ascending {
		return p.Fitness(p.items[i]) < p.Fitness(p.items[j])
	}
	return p.Fitness(p.items[i]) > p.Fitness(p.items[j])
}

// heap.Interface
func (p *Priority) Push(x interface{}) {
	if p.length+1 > len(p.items) {
		tmp := make([]*Pattern, len(p.items)+50000)
		copy(tmp, p.items)
		p.items = tmp
	}

	p.items[p.length] = x.(*Pattern)
	p.length += 1
}

func (p *Priority) Pop() interface{} {
	r := p.items[p.length-1]
	p.length -= 1
	return r
}
