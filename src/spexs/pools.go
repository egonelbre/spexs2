package spexs

import (
	"container/heap"
	"container/list"
)

type FifoPool struct {
	token chan int
	list  *list.List
}

func NewFifoPool() *FifoPool {
	p := &FifoPool{}
	p.token = make(chan int, 1)
	p.list = list.New()
	p.token <- 1
	return p
}

func (p *FifoPool) Take() (*TrieNode, bool) {
	<-p.token
	if p.list.Len() == 0 {
		p.token <- 1
		return nil, false
	}
	tmp := p.list.Front()
	p.list.Remove(tmp)
	p.token <- 1
	return tmp.Value.(*TrieNode), true
}

func (p *FifoPool) Put(pat *TrieNode) {
	<-p.token
	p.list.PushBack(pat)
	p.token <- 1
}

func (p *FifoPool) Len() int {
	return p.list.Len()
}

type TrieFitnessFunc func(p *TrieNode) float64

type PriorityPool struct {
	token   chan int
	items   []*TrieNode
	Fitness TrieFitnessFunc
	limit   int
}

func NewPriorityPool(fitness TrieFitnessFunc, limit int) *PriorityPool {
	p := &PriorityPool{}
	p.token = make(chan int, 1)
	p.items = make([]*TrieNode, 0)
	p.limit = limit
	p.Fitness = fitness
	p.token <- 1

	heap.Init(p)
	return p
}

func (p *PriorityPool) IsEmpty() bool {
	return len(p.items) == 0
}

func (p *PriorityPool) Take() (*TrieNode, bool) {
	<-p.token
	if p.IsEmpty() {
		p.token <- 1
		return nil, false
	}
	v := heap.Pop(p)
	p.token <- 1
	return v.(*TrieNode), true
}

func (p *PriorityPool) Put(pat *TrieNode) {
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
	p.items = append(p.items, x.(*TrieNode))
}

func (p *PriorityPool) Pop() interface{} {
	r := p.items[len(p.items)-1]
	p.items = p.items[0 : len(p.items)-1]
	return r
}
