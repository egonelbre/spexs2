package trie

const (
	patternsBufferSize = 128
)

type Char rune

type Patterns chan *Pattern

type Pooler interface {
	Take() (*Pattern, bool)
	Put(*Pattern)
	Len() int
}

type ExtenderFunc func(p *Pattern, ref *Reference) Patterns
type FilterFunc func(p *Pattern, ref *Reference) bool
type FitnessFunc func(p *Pattern) float64

type Setup struct {
	Ref *Reference
	Out Pooler
	In  Pooler

	Extender ExtenderFunc

	Extendable  FilterFunc
	Outputtable FilterFunc
}

func NewPatterns() Patterns {
	return make(Patterns, patternsBufferSize)
}

func Run(s Setup) {
	for {
		p, valid := s.In.Take()
		if !valid {
			return
		}

		extensions := s.Extender(p, s.Ref)
		for extended := range extensions {
			if s.Extendable(extended, s.Ref) {
				s.In.Put(extended)
			}
			if s.Outputtable(extended, s.Ref) {
				s.Out.Put(extended)
			}
		}
	}
}

func Parallel(f func(), routines int){
	stop := make(chan int, routines)
	for i := 0; i < routines; i += 1 {
		go func() {
			defer func() { stop <- 1 }()
			f()
		}()
	}
	for i := 0; i < routines; i += 1 {
		<-stop
	}
}

func RunParallel(s Setup, routines int) {
	Parallel(func(){Run(s)}, routines)
}
