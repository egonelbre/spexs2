package spexs

type Pattern interface {
	ToString() string
}

type Patterns chan *Pattern

var EmptyPattern *Pattern = nil

type Reference interface {
	Next(pos Pos) (Char, Pos, bool)
}

type Pooler interface {
	Take() (*Pattern, bool)
	Put(*Pattern)
}

type PatternFilter func(p *Pattern) bool
type ExtenderFunc func(p *Pattern, ref *Reference) Patterns

func Run(ref Reference, patterns Pooler, results Pooler,
	extender ExtenderFunc, acceptable PatternFilter) {
	patterns.Put(EmptyPattern)
	p, valid := patterns.Take()
	for valid {
		pats := extender(p, &ref)
		for ep := range pats {
			if acceptable(ep) {
				patterns.Put(ep)
				results.Put(ep)
			}
		}
		p, valid = patterns.Take()
	}
}

func RunParallel(ref Reference, input Pooler, results Pooler,
	extender ExtenderFunc, acceptable PatternFilter) {
	input.Put(EmptyPattern)
	start := make(chan int) // alternatively make(chan int, threadLimit)
	stop := make(chan int)
	
	num_threads := 10;
	for i := 0; i < num_threads; i++ {
		go func(){
			start <- 1
			defer func() { stop <- 1 }()

			for {
				p, valid := input.Take()
				if !valid {
					return
				}

				pats := extender(p, &ref)
				for ep := range pats {
					if acceptable(ep) {
						input.Put(ep)
						results.Put(ep)
					}
				}
			}
		}()
	}

	for _ = range(start) {
		<-stop
	}
}
