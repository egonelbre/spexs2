package spexs

import "sync"

type Token uint32
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

type signal struct{}

func RunParallel(s *Setup, routines int) {
	prepareSpexs(s)

	wg := sync.WaitGroup{}

	allDone := false
	m, out := &sync.Mutex{}, &sync.Mutex{}

	added := make(chan signal, 1e9)
	added <- signal{}
	workers := 0

	for i := 0; i < routines; i += 1 {
		wg.Add(1)

		go func() {
			for {
				<-added
				m.Lock()
				if allDone {
					added <- signal{}
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

						added <- signal{}

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
					added <- signal{}
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
