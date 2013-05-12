package utils

import "sync"

type Sem struct {
	mutex sync.Mutex
	delay sync.Mutex
	count int32
}

func NewSem(n int32) *Sem {
	s := &Sem{}
	s.delay.Lock()
	s.count = n
	return s
}

func (s *Sem) Wait() {
	s.mutex.Lock()
	s.count -= 1
	if s.count < 0 {
		s.mutex.Unlock()
		s.delay.Lock()
	}
	s.mutex.Unlock()
	return
}

func (s *Sem) Signal() {
	s.mutex.Lock()
	s.count += 1
	if s.count <= 0 {
		s.delay.Unlock()
	} else {
		s.mutex.Unlock()
	}
}
