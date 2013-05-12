package utils

type signal struct{}

type Sem struct {
	some  chan signal
	none  chan signal
	count uint32
}

func NewSem(n uint32) *Sem {
	s := Sem{make(chan signal, 1), make(chan signal, 1), n}
	if n == 0 {
		s.none <- signal{}
	} else {
		s.some <- signal{}
	}
	return &s
}

func (s *Sem) Wait() {
	<-s.some
	s.count--
	if s.count == 0 {
		s.none <- signal{}
	} else {
		s.some <- signal{}
	}
}

func (s *Sem) Signal() {
	select {
	case <-s.some:
	case <-s.none:
	}
	s.count++
	s.some <- signal{}
}
