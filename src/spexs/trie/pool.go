package trie

import (
	"container/heap"
)

type PriorityPool struct {
	token     chan int
	items     []*Pattern
	Fitness   FitnessFunc
	limit     int
	ascending bool
}

func NewPriorityPool(fitness FitnessFunc, limit int, ascending bool) *PriorityPool {
	p := &PriorityPool{}
	p.token = make(chan int, 1)
	p.items = make([]*Pattern, 0)
	p.limit = limit
	p.Fitness = fitness
	p.token <- 1
	p.ascending = ascending

	heap.Init(p)
	return p
}

func (p *PriorityPool) IsEmpty() bool {
	return len(p.items) == 0
}

func (p *PriorityPool) Take() (*Pattern, bool) {
	<-p.token
	if p.IsEmpty() {
		p.token <- 1
		return nil, false
	}
	v := heap.Pop(p)
	p.token <- 1
	return v.(*Pattern), true
}

func (p *PriorityPool) Put(pat *Pattern) {
	<-p.token
	p.token <- 1
	return
	heap.Push(p, pat)
	if p.limit > 0 && p.Len() > p.limit {
		heap.Pop(p)
	}
	p.token <- 1
}

// sort.Interface
func (p *PriorityPool) Len() int {
	return len(p.items)
}

func (p *PriorityPool) Swap(i, j int) {
	temp := p.items[i]
	p.items[i] = p.items[j]
	p.items[j] = temp
}

func (p *PriorityPool) Less(i, j int) bool {
	if p.ascending {
		return p.Fitness(p.items[i]) > p.Fitness(p.items[j])
	}
	return p.Fitness(p.items[i]) < p.Fitness(p.items[j])
}

// heap.Interface
func (p *PriorityPool) Push(x interface{}) {
	p.items = append(p.items, x.(*Pattern))
}

func (p *PriorityPool) Pop() interface{} {
	r := p.items[len(p.items)-1]
	p.items = p.items[0 : len(p.items)-1]
	return r
}
