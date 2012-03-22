package spexs

import "fmt"

const (
	patternsBufferSize = 1024*1024*204
)

type Pattern interface {
	fmt.Stringer
}

type Patterns chan Pattern

func MakePatterns() Patterns {
	return make(Patterns, patternsBufferSize)
}

type Reference interface {
	Next(pos Pos) (Char, Pos, bool)
}

type Pooler interface {
	Take() (Pattern, bool)
	Put(Pattern)
}

type PatternFilter func(p Pattern) bool
type ExtenderFunc func(p Pattern, ref Reference) Patterns

func Run(ref Reference, patterns Pooler, results Pooler,
	extender ExtenderFunc, acceptable PatternFilter) {
	p, valid := patterns.Take()
	for valid {
		pats := extender(p, ref)
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

	start := make(chan int, 1000)
	stop := make(chan int, 1000)
	
	num_threads := 36;
	for i := 0; i < num_threads; i++ {
		go func(){
			start <- 1
			defer func() { stop <- 1 }()

			for {
				p, valid := input.Take()
				if !valid {
					return
				}

				pats := extender(p, ref)
				for ep := range pats {
					if acceptable(ep) {
						input.Put(ep)
						results.Put(ep)
					}
				}
			}
		}()
	}

	for i := 0; i < num_threads; i++ {
		<-stop
	}
}
