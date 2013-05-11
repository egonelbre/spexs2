package spexs

import "time"

type Token uint
type Querys []*Query

type Pooler interface {
	Pop() (*Query, bool)
	Push(*Query)
	Values() []*Query
	Len() int
}

type Extender func(p *Query) Querys
type Filter func(p *Query) bool
type ProcessQuery func(p *Query) error
type Feature func(p *Query) (float64, string)

type Setup struct {
	Db  *Database
	Out Pooler
	In  Pooler

	Extender Extender

	Extendable  Filter
	Outputtable Filter

	PreProcess  ProcessQuery
	PostProcess ProcessQuery
}

func prepareSpexs(s *Setup) {
	maxSeq := 0
	for _, seq := range s.Db.Sequences {
		length := len(seq.Tokens)
		if length > maxSeq {
			maxSeq = length
		}
	}

	for i := uint(0); i < 32; i += 1 {
		if 1<<i > maxSeq {
			PosOffset = i
			break
		}
	}

	s.In.Push(NewEmptyQuery(s.Db))
}

func Run(s *Setup) {
	prepareSpexs(s)
	for {
		p, ok := s.In.Pop()
		if !ok {
			return
		}

		extensions := s.Extender(p)
		for _, extended := range extensions {
			if s.Extendable(extended) {
				s.In.Push(extended)
				if s.Outputtable(extended) {
					s.Out.Push(extended)
				}
			}
		}

		if s.PostProcess(p) != nil {
			break
		}
	}
}

func RunParallel(s *Setup, routines int) {
	prepareSpexs(s)

	quit := make(chan int, routines)
	counter := make(chan int, routines)

	for i := 0; i < routines; i += 1 {
		go func(rtn int) {
		main:
			for {
				p, ok := s.In.Pop()
				for !ok {
					p, ok = s.In.Pop()
					if ok {
						break
					}
					counter <- -1
					select {
					case <-time.After(100 * time.Millisecond):
					case <-quit:
						break main
					}
					counter <- 1
				}

				extensions := s.Extender(p)
				for _, extended := range extensions {
					if s.Extendable(extended) {
						s.PreProcess(extended)
						s.In.Push(extended)
						if s.Outputtable(extended) {
							s.Out.Push(extended)
						}
					}
				}
				if s.PostProcess(p) != nil {
					break
				}
			}
		}(i)
	}

	count := routines
	for count > 0 {
		value := <-counter
		count += value
	}

	for i := 0; i < routines; i += 1 {
		quit <- 1
	}
}
