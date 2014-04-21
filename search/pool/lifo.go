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

// Package stack implements a LIFO (last in first out) data structure supporting
// arbitrary types (even a mixture).
//
// Internally it uses a dynamically growing slice of blocks, resulting in faster
// resizes than a simple dynamic array/slice would allow.
package pool

import (
	"github.com/egonelbre/spexs2/search"
)

// Last in, first out data structure.
type Stack struct {
	size     int
	capacity int
	offset   int

	blocks [][]*search.Query
	active []*search.Query
}

// Creates a new, empty stack.
func NewStack() *Stack {
	result := new(Stack)
	result.active = make([]*search.Query, blockSize)
	result.blocks = [][]*search.Query{result.active}
	result.capacity = blockSize
	return result
}

// Pushes a value onto the stack, expanding it if necessary.
func (s *Stack) Push(data *search.Query) {
	if s.size == s.capacity {
		s.active = make([]*search.Query, blockSize)
		s.blocks = append(s.blocks, s.active)
		s.capacity += blockSize
		s.offset = 0
	} else if s.offset == blockSize {
		s.active = s.blocks[s.size/blockSize]
		s.offset = 0
	}
	s.active[s.offset] = data
	s.offset++
	s.size++
}

// Pops a value off the stack and returns it. Currently no shrinking is done.
func (s *Stack) Pop() (res *search.Query, ok bool) {
	if s.size == 0 {
		return nil, false
	}
	s.size--
	s.offset--
	if s.offset < 0 {
		s.offset = blockSize - 1
		s.active = s.blocks[s.size/blockSize]
	}
	res, s.active[s.offset] = s.active[s.offset], nil
	return res, true
}

// Checks whether the stack is empty or not.
func (s *Stack) Empty() bool {
	r := s.size == 0
	return r
}

// Returns all values in an array
func (s *Stack) Values() []*search.Query {
	r := make([]*search.Query, 0, s.Len())
	for !s.Empty() {
		v, _ := s.Pop()
		r = append(r, v)
	}
	return r
}

// Returns the number of elements in the stack.
func (s *Stack) Len() int {
	r := s.size
	return r
}
