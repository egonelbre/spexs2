package search

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

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
type ProcessQuery func(p *Query)
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

func Run(s *Setup) {
	s.In.Push(NewEmptyQuery(s.Db))
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

		s.PostProcess(p)
	}
}

type signal struct{}

func RunParallel(s *Setup, routines int) {
	s.In.Push(NewEmptyQuery(s.Db))

	wg := sync.WaitGroup{}

	allDone := false
	m, out := &sync.Mutex{}, &sync.Mutex{}

	added := make(chan signal, 1e9)
	added <- signal{}
	workers := 0

	for i := 0; i < routines; i += 1 {
		wg.Add(1)
		go func() {
			runtime.LockOSThread()
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
				s.PostProcess(p)

				if allDone {
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

func pipe(from chan *Query, storage Pooler, to chan *Query) {
	var next *Query
	for {
		switch {
		case next != nil:
			select {
			case q, ok := <-from:
				if !ok {
					return
				}
				storage.Push(q)
			case to <- next:
				next = nil
			}
		case storage.Len() > 0:
			next, _ = storage.Pop()
		default:
			q, ok := <-from
			if !ok {
				return
			}
			next = q
		}
	}
}

func RunParallelChan(s *Setup, procs int) {
	const Buffer = 128

	// work queue for workers
	work := make(chan *Query, procs)
	// this pumps elements into the input pool
	newwork := make(chan *Query, Buffer)
	// this pumps results into output pool
	found := make(chan *Query, Buffer)

	// sort the work items
	go pipe(newwork, s.In, work)

	// monitor found queries
	var output sync.WaitGroup
	output.Add(1)
	go func() {
		for q := range found {
			s.Out.Push(q)
		}
		output.Done()
	}()

	// number of unfinished items
	var pending int32 = 1
	work <- NewEmptyQuery(s.Db)

	for i := 0; i < procs; i += 1 {
		go func() {
			for p := range work {
				extensions := s.Extender(p)
				for _, extended := range extensions {
					if s.Extendable(extended) {
						s.PreProcess(extended)

						atomic.AddInt32(&pending, 1)
						newwork <- extended

						if s.Outputtable(extended) {
							found <- extended
						}
					}
				}
				s.PostProcess(p)
				atomic.AddInt32(&pending, -1)
			}
		}()
	}

	// monitor whether there is no work pending
	for _ = range time.Tick(10 * time.Millisecond) {
		if atomic.LoadInt32(&pending) == 0 {
			break
		}
	}

	close(work)
	close(newwork)
	close(found)

	output.Wait()
}
