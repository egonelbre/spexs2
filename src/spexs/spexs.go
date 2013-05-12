package spexs

import (
	"sync"
	"utils"
)

type Token uint
type Querys []*Query

type Pooler interface {
	Pop() (*Query, bool)
	Push(*Query)
	Values() []*Query
	Len() int
	Empty() bool
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

	wg := sync.WaitGroup{}

	allDone := false
	m, out := &sync.Mutex{}, &sync.Mutex{}

	added := utils.NewSem(1)
	workers := 0

	for i := 0; i < routines; i += 1 {
		wg.Add(1)

		go func() {
			for {
				added.Wait()
				m.Lock()
				if allDone {
					added.Signal()
					m.Unlock()
					break
				}

				p, ok := s.In.Pop()
				if !ok {
					m.Unlock()
					continue
				}
				workers += 1
				m.Unlock()

				extensions := s.Extender(p)
				for _, extended := range extensions {
					if s.Extendable(extended) {
						s.PreProcess(extended)

						m.Lock()
						s.In.Push(extended)
						m.Unlock()

						added.Signal()

						if s.Outputtable(extended) {
							out.Lock()
							s.Out.Push(extended)
							out.Unlock()
						}
					}
				}

				m.Lock()
				workers -= 1
				allDone = workers == 0 && s.In.Empty()
				needToTerminate := s.PostProcess(p) != nil

				if allDone || needToTerminate {
					added.Signal()
					m.Unlock()
					break
				}
				m.Unlock()
			}

			wg.Done()
		}()
	}

	wg.Wait()
}
