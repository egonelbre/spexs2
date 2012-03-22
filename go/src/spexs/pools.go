package spexs

import "container/heap"

type FifoPool struct {
	patterns Patterns
}

func NewFifoPool() *FifoPool {
	return &FifoPool{MakePatterns()}
}

func (p *FifoPool) Take() (Pattern, bool) {
	select {
		case pat := <-p.patterns:
			return pat, true
		default: break
	}
	return nil, false
}

func (p *FifoPool) Put(pat Pattern) {
	p.patterns <- pat
}

type FitnessFunc func(a Pattern) float32

type PriorityPool struct {
	token   chan int
	items   []Pattern
	Fitness FitnessFunc
	limit   int
}

func NewPriorityPool(fitness FitnessFunc, limit int) *PriorityPool {
	p := &PriorityPool{}
	p.token = make(chan int, 1)
	p.items = make([]Pattern, 0)
	p.limit = limit
	p.Fitness = fitness
	p.token <- 1

	heap.Init(p)
	return p
}

func (p *PriorityPool) IsEmpty() bool {
    return len(p.items) == 0
}

func (p *PriorityPool) Take() (Pattern, bool) {
	<-p.token
	defer func() { p.token <- 1 }()

	if p.IsEmpty() {
		return nil, false
	}
	v := heap.Pop(p)
	return v.(Pattern), true
}

func (p *PriorityPool) Put(pat Pattern) {
	<-p.token
	defer func() { p.token <- 1 }()

	heap.Push(p, pat)
	if p.Len() > p.limit {
		heap.Pop(p)
	}
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
	return p.Fitness(p.items[i]) < p.Fitness(p.items[j])
}

// heap.Interface
func (p *PriorityPool) Push(x interface{}) {
	p.items = append(p.items, x.(Pattern))
}

func (p *PriorityPool) Pop() interface{} {
	r := p.items[len(p.items)-1]
	p.items = p.items[0 : len(p.items)-1]
	return r
}
