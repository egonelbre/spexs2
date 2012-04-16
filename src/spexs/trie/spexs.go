package trie

type Patterns chan Pattern

type Pooler interface{}{
	Take() (Pattern, bool)
	Put(Pattern)
	Len() int
}

type FilterFunc func(p Pattern) bool
type ExtenderFunc func(p Pattern, ref Reference) Patterns

type Setup struct {
	Ref Reference
	Out Pooler
	In  Pooler

	Extender ExtenderFunc

	Extendable  FilterFunc
	Outputtable FilterFunc
}

func Run(s setup){
	for {
		p, valid := s.In.Take()
		if !valid {
			return
		}

		extensions := s.Extender(p, s.Ref)
		for extension := range extensions {
			if s.Extendable(extension) {
				s.In.Put(extension)
			}
			if s.Outputtable(extension) {
				s.Out.Put(extension)
			}
		}
	}
}

func RunParallel(s setup, routines int){
	stop := make(chan int, routines)

	for i := 0; i < routines; i += 1 {
		go func() {
			defer func() { stop <- 1}
			Run(s)
		}()
	}

	for i := 0; i < routines; i += 1 {
		<-stop
	}
}