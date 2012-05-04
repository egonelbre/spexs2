package pool

import (
	"spexs"
	"sync/atomic"
	"unsafe"
)

type node struct {
	val *spexs.Pattern
	nxt unsafe.Pointer
}

type Lifo struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

func NewLifo() (q *Lifo) {
	q = new(Lifo)
	n := unsafe.Pointer(new(node))
	q.head = n
	q.tail = n
	return
}

func (q *Lifo) Take() (val *spexs.Pattern, success bool) {
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

func (q *Lifo) Put(val *spexs.Pattern) {
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

func (p *Lifo) Len() int {
	return 1
}
