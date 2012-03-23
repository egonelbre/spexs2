package spexs

import (
	"container/list"
	"container/heap"
)

type FifoPool struct {
	token chan int
	list *list.List
}

func NewFifoPool() *FifoPool {
	p := &FifoPool{}
	p.token = make(chan int, 1)
	p.list = list.New()
	p.token <- 1
	return p
}

func (p *FifoPool) Take() (Pattern, bool) {
	<-p.token
	if p.list.Len() == 0 {
		p.token <- 1
		return nil, false
	}
	tmp := p.list.Front()
	p.list.Remove(tmp)
	p.token <- 1
	return tmp.Value.(Pattern), true
}

func (p *FifoPool) Put(pat Pattern) {
	<-p.token
	p.list.PushBack(pat)
	p.token <- 1
}

func (p *FifoPool) Len() int {
	return p.list.Len()
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
	if p.IsEmpty() {
		p.token <- 1
		return nil, false
	}
	v := heap.Pop(p)
	p.token <- 1
	return v.(Pattern), true
}

func (p *PriorityPool) Put(pat Pattern) {
	<-p.token
	heap.Push(p, pat)
	if p.Len() > p.limit {
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
