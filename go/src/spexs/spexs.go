package spexs

import "fmt"

const (
	patternsBufferSize = 100
)

type Pattern interface {
	fmt.Stringer
}

type Patterns chan Pattern

func MakePatterns() Patterns {
	return make(Patterns, patternsBufferSize)
}

type Reference interface {
	Next(idx int, pos byte) (Char, byte, bool)
}

type Pooler interface {
	Take() (Pattern, bool)
	Put(Pattern)
	Len() int
}

type FilterFunc func(p Pattern) bool
type ExtenderFunc func(p Pattern, ref Reference) Patterns

func Run(ref Reference, input Pooler, results Pooler,
	extender ExtenderFunc, extendable FilterFunc, outputtable FilterFunc) {
	p, valid := input.Take()
	for valid {
		pats := extender(p, ref)
		for ep := range pats {
			if extendable(ep) {
				input.Put(ep)
			}
			if outputtable(ep) {
				results.Put(ep)
			}
		}
		p, valid = input.Take()
	}
}

func RunParallel(ref Reference, input Pooler, results Pooler,
	extender ExtenderFunc, extendable FilterFunc, outputtable FilterFunc, num_threads int) {

	start := make(chan int, 1000)
	stop := make(chan int, 1000)

	for i := 0; i < num_threads; i++ {
		go func() {
			start <- 1
			defer func() { stop <- 1 }()

			for {
				p, valid := input.Take()
				if !valid {
					return
				}

				pats := extender(p, ref)
				for ep := range pats {
					if extendable(ep) {
						input.Put(ep)
					}
					if outputtable(ep) {
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
