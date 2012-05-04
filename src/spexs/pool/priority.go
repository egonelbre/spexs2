package pool

import (
	"spexs"
	"container/heap"
)

type Priority struct {
	token     chan int
	items     []*spexs.Pattern
	Fitness   spexs.FitnessFunc
	limit     int
	ascending bool
}

func NewPriority(fitness spexs.FitnessFunc, limit int, ascending bool) *Priority {
	p := &Priority{}
	p.token = make(chan int, 1)
	p.items = make([]*spexs.Pattern, 0)
	p.limit = limit
	p.Fitness = fitness
	p.token <- 1
	p.ascending = ascending

	heap.Init(p)
	return p
}

func (p *Priority) IsEmpty() bool {
	return len(p.items) == 0
}

func (p *Priority) Take() (*spexs.Pattern, bool) {
	<-p.token
	if p.IsEmpty() {
		p.token <- 1
		return nil, false
	}
	v := heap.Pop(p)
	p.token <- 1
	return v.(*spexs.Pattern), true
}

func (p *Priority) Put(pat *spexs.Pattern) {
	<-p.token
	heap.Push(p, pat)
	if p.limit > 0 && p.Len() > p.limit {
		heap.Pop(p)
	}
	p.token <- 1
}

// sort.Interface
func (p *Priority) Len() int {
	return len(p.items)
}

func (p *Priority) Swap(i, j int) {
	temp := p.items[i]
	p.items[i] = p.items[j]
	p.items[j] = temp
}

func (p *Priority) Less(i, j int) bool {
	if p.ascending {
		return p.Fitness(p.items[i]) < p.Fitness(p.items[j])
	}
	return p.Fitness(p.items[i]) > p.Fitness(p.items[j])
}

// heap.Interface
func (p *Priority) Push(x interface{}) {
	p.items = append(p.items, x.(*spexs.Pattern))
}

func (p *Priority) Pop() interface{} {
	r := p.items[len(p.items)-1]
	p.items = p.items[0 : len(p.items)-1]
	return r
}
