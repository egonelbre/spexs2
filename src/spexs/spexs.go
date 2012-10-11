package spexs

import "time"

type Token uint
type Querys []*Query

type Pooler interface {
	Take() (*Query, bool)
	Put(*Query)
	Len() int
}

type ExtenderFunc func(p *Query, db *Database) Querys
type FilterFunc func(p *Query, db *Database) bool
type FitnessFunc func(p *Query) float64
type PostProcessFunc func(p *Query, s *Setup) error

type Setup struct {
	DB  *Database
	Out Pooler
	In  Pooler

	Extender ExtenderFunc

	Extendable  FilterFunc
	Outputtable FilterFunc

	PostProcess PostProcessFunc
}

func prepareSpexs(s *Setup) {
	maxSeq := 0
	for _, seq := range s.DB.Sequences {
		if seq.Len > maxSeq {
			maxSeq = seq.Len
		}
	}

	for i := uint(0); i < 32; i += 1 {
		if 1<<i > maxSeq {
			PosOffset = i
			break
		}
	}

	s.In.Put(NewEmptyQuery(s.DB))
}

func Run(s *Setup) {
	prepareSpexs(s)
	for {
		p, valid := s.In.Take()
		if !valid {
			return
		}

		extensions := s.Extender(p, s.DB)
		for _, extended := range extensions {
			if s.Extendable(extended, s.DB) {
				s.In.Put(extended)
			}
			if s.Outputtable(extended, s.DB) {
				s.Out.Put(extended)
			}
		}

		if s.PostProcess(p, s) != nil {
			break
		}
	}
}

func Parallel(f func(), routines int) {
	stop := make(chan int, routines)
	for i := 0; i < routines; i += 1 {
		go func(rtn int) {
			defer func() { stop <- 1 }()
			f()
		}(i)
	}
	for i := 0; i < routines; i += 1 {
		<-stop
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
				p, valid := s.In.Take()
				for !valid {
					p, valid = s.In.Take()
					if valid {
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

				extensions := s.Extender(p, s.DB)
				for _, extended := range extensions {
					if s.Extendable(extended, s.DB) {
						s.In.Put(extended)
					}
					if s.Outputtable(extended, s.DB) {
						s.Out.Put(extended)
					}
				}

				if s.PostProcess(p, s) != nil {
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
