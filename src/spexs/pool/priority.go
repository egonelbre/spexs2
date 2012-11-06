package pool

import (
	"container/heap"
	"sort"
	. "spexs"
)

type Priority struct {
	token  chan int
	items  []*Query
	Order  Feature
	length int
	limit  int
}

func NewPriority(order Feature, limit int) *Priority {
	p := &Priority{}
	p.token = make(chan int, 1)
	p.items = make([]*Query, limit+100)
	p.length = 0
	p.limit = limit
	p.Order = order
	p.token <- 1

	heap.Init(p)
	return p
}

func (p *Priority) IsEmpty() bool {
	return p.length == 0
}

func (p *Priority) Take() (*Query, bool) {
	<-p.token
	if p.IsEmpty() {
		p.token <- 1
		return nil, false
	}
	v := heap.Pop(p)
	p.token <- 1
	return v.(*Query), true
}

func (p *Priority) Put(pat *Query) {
	<-p.token
	heap.Push(p, pat)
	//if p.limit > 0 && p.Len() > p.limit {
	//	heap.Pop(p)
	//}
	p.token <- 1
}

func (p *Priority) Top(n int) []*Query {
	sort.Sort(p)
	last := n
	if last > p.length {
		last = p.length
	}
	return p.items[:last]
}

func (p *Priority) Bottom(n int) []*Query {
	sort.Sort(p)
	first := p.length - n
	if first < 0 {
		first = 0
	}
	return p.items[first:p.length]
}

func (p *Priority) Heapify() {
	heap.Init(p)
}

// sort.Interface
func (p *Priority) Len() int {
	return p.length
}

func (p *Priority) Swap(i, j int) {
	p.items[i], p.items[j] = p.items[j], p.items[i]
}

func (p *Priority) Less(i, j int) bool {
	a, _ := p.items[i].Memoized(p.Order)
	b, _ := p.items[j].Memoized(p.Order)
	return a < b
}

// heap.Interface
func (p *Priority) Push(x interface{}) {
	if p.length+1 > len(p.items) {
		tmp := make([]*Query, len(p.items)+50000)
		copy(tmp, p.items)
		p.items = tmp
	}

	p.items[p.length] = x.(*Query)
	p.length += 1
}

func (p *Priority) Pop() interface{} {
	r := p.items[p.length-1]
	p.length -= 1
	return r
}
