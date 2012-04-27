package trie

import (
	"container/list"
	"sync/atomic"
	"unsafe"
)

type node struct {
	val *Pattern
	nxt unsafe.Pointer
}

type LifoPool struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

func NewLifoPool() (q *LifoPool) {
	q = new(LifoPool)
	n := unsafe.Pointer(new(node))
	q.head = n
	q.tail = n
	return
}

func (q *LifoPool) Take() (val *Pattern, success bool) {
	var h, t, n unsafe.Pointer
	for {
		h = q.head
		t = q.tail
		n = ((*node)(h)).nxt
		if h == t {
			if n == nil {
				return nil, false
			} else {
				atomic.CompareAndSwapPointer(&q.tail, t, n)
			}
		} else {
			val = ((*node)(n)).val // Enq(...) write to val may not be visible
			if atomic.CompareAndSwapPointer(&q.head, h, n) {
				return val, true
			}
		}
	}
	panic("Unreachable")
}

func (q *LifoPool) Put(val *Pattern) {
	var t, n unsafe.Pointer
	n = unsafe.Pointer(&node{val: val, nxt: nil})
	for {
		t = q.tail
		nxt := ((*node)(t)).nxt
		if nxt != nil {
			atomic.CompareAndSwapPointer(&q.tail, t, nxt)
		} else if atomic.CompareAndSwapPointer(&((*node)(t)).nxt, nil, n) {
			break
		}
	}
	atomic.CompareAndSwapPointer(&q.tail, t, n)
}

func (p *LifoPool) Len() int {
	return 1
}

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

func (p *FifoPool) Take() (*Pattern, bool) {
	<-p.token
	if p.list.Len() == 0 {
		p.token <- 1
		return nil, false
	}
	tmp := p.list.Front()
	p.list.Remove(tmp)
	p.token <- 1
	return tmp.Value.(*Pattern), true
}

func (p *FifoPool) Put(pat *Pattern) {
	<-p.token
	p.list.PushBack(pat)
	p.token <- 1
}

func (p *FifoPool) Len() int {
	return p.list.Len()
}
