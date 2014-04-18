// CookieJar - A contestant's algorithm toolbox
// Copyright (c) 2013 Peter Szilagyi. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
//     * Redistributions of source code must retain the above copyright notice,
//       this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above copyright
//       notice, this list of conditions and the following disclaimer in the
//       documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.
//
// Alternatively, the CookieJar toolbox may be used in accordance with the terms
// and conditions contained in a signed written agreement between you and the
// author(s).
//
// Author: peterke@gmail.com (Peter Szilagyi)

// Package queue implements a FIFO (first in first out) data structure supporting
// arbitrary types (even a mixture).
//
// Internally it uses a dynamically growing circular slice of blocks, resulting
// in faster resizes than a simple dynamic array/slice would allow.
package pool

import (
	"github.com/egonelbre/spexs2/search"
)

// First in, first out data structure.
type Queue struct {
	tailIdx int
	headIdx int
	tailOff int
	headOff int

	blocks [][]*search.Query
	head   []*search.Query
	tail   []*search.Query
}

// Creates a new, empty queue.
func NewQueue() *Queue {
	result := new(Queue)
	result.blocks = [][]*search.Query{make([]*search.Query, blockSize)}
	result.head = result.blocks[0]
	result.tail = result.blocks[0]
	return result
}

// Pushes a new element into the queue, expanding it if necessary.
func (q *Queue) Push(data *search.Query) {
	q.tail[q.tailOff] = data
	q.tailOff++
	if q.tailOff == blockSize {
		q.tailOff = 0
		q.tailIdx = (q.tailIdx + 1) % len(q.blocks)

		// If we wrapped over to the end, insert a new block and update indices
		if q.tailIdx == q.headIdx {
			buffer := make([][]*search.Query, len(q.blocks)+1)
			copy(buffer[:q.tailIdx], q.blocks[:q.tailIdx])
			buffer[q.tailIdx] = make([]*search.Query, blockSize)
			copy(buffer[q.tailIdx+1:], q.blocks[q.tailIdx:])
			q.blocks = buffer
			q.headIdx++
			q.head = q.blocks[q.headIdx]
		}
		q.tail = q.blocks[q.tailIdx]
	}
}

// Pops out an element from the queue. Note, no bounds checking are done.
func (q *Queue) Pop() (res *search.Query, ok bool) {
	if q.headIdx == q.tailIdx && q.headOff == q.tailOff {
		return nil, false
	}
	res, q.head[q.headOff] = q.head[q.headOff], nil
	q.headOff++
	if q.headOff == blockSize {
		q.headOff = 0
		q.headIdx = (q.headIdx + 1) % len(q.blocks)
		q.head = q.blocks[q.headIdx]
	}
	return res, true
}

// Checks whether the queue is empty.
func (q *Queue) Empty() bool {
	r := q.headIdx == q.tailIdx && q.headOff == q.tailOff
	return r
}

// Returns the number of elements in the queue.
func (q *Queue) Len() (size int) {
	if q.tailIdx > q.headIdx {
		size = (q.tailIdx-q.headIdx)*blockSize - q.headOff + q.tailOff
	} else if q.tailIdx < q.headIdx {
		size = (len(q.blocks)-q.headIdx+q.tailIdx)*blockSize - q.headOff + q.tailOff
	} else {
		size = q.tailOff - q.headOff
	}
	return size
}

// Returns all values in an array
func (q *Queue) Values() []*search.Query {
	r := make([]*search.Query, 0, q.Len())
	for !q.Empty() {
		v, _ := q.Pop()
		r = append(r, v)
	}
	return r
}
