package pool

import (
	"container/heap"
	"sort"

	"github.com/egonelbre/spexs2/search"
)

type Priority struct {
	items  []*search.Query
	Order  []search.Feature
	Worst  *search.Query
	length int
	limit  int
}

type priorityIntf Priority

func NewPriority(order []search.Feature, limit int) *Priority {
	p := &Priority{}
	p.items = make([]*search.Query, limit+100)
	p.length = 0
	p.limit = limit
	p.Order = order

	heap.Init((*priorityIntf)(p))
	return p
}

func (p *Priority) Empty() bool {
	return p.length == 0
}

func (p *Priority) Pop() (*search.Query, bool) {
	if p.Empty() {
		return nil, false
	}
	v := heap.Pop((*priorityIntf)(p))
	return v.(*search.Query), true
}

func (p *Priority) Push(pat *search.Query) {
	worst := p.Worst
	if worst != nil && p.less(pat, worst) {
		return
	}

	heap.Push((*priorityIntf)(p), pat)
	if p.limit > 0 && p.Len() > p.limit {
		p.Worst = heap.Pop((*priorityIntf)(p)).(*search.Query)
	}
}

func (p *Priority) Top(n int) []*search.Query {
	sort.Sort((*priorityIntf)(p))
	last := n
	if last > p.length {
		last = p.length
	}
	return p.items[:last]
}

func (p *Priority) Bottom(n int) []*search.Query {
	sort.Sort((*priorityIntf)(p))
	first := p.length - n
	if first < 0 {
		first = 0
	}
	items := p.items[first:p.length]
	n = len(items)
	result := make([]*search.Query, n)
	for i := 0; i < n; i++ {
		result[i] = items[n-i-1]
	}
	return result
}

func (p *Priority) Values() []*search.Query {
	return p.Bottom(p.limit)
}

func (p *Priority) Heapify() {
	heap.Init((*priorityIntf)(p))
}

func (p *Priority) less(a *search.Query, b *search.Query) bool {
	for _, fn := range p.Order {
		aval, _ := fn(a)
		bval, _ := fn(b)

		if aval != bval {
			return aval < bval
		}
	}
	return false
}

// sort.Interface
func (p *Priority) Len() int {
	return p.length
}

func (p *priorityIntf) Len() int {
	return p.length
}

func (p *priorityIntf) Swap(i, j int) {
	p.items[i], p.items[j] = p.items[j], p.items[i]
}

func (p *priorityIntf) Less(i, j int) bool {
	return (*Priority)(p).less(p.items[i], p.items[j])
}

// heap.Interface
func (p *priorityIntf) Push(x interface{}) {
	if p.length+1 > len(p.items) {
		tmp := make([]*search.Query, len(p.items)+50000)
		copy(tmp, p.items)
		p.items = tmp
	}

	p.items[p.length] = x.(*search.Query)
	p.length++
}

func (p *priorityIntf) Pop() interface{} {
	r := p.items[p.length-1]
	p.length--
	return r
}
