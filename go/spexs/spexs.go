package spexs

type Pattern interface{
  ToString() string
}

type Patterns chan Pattern
type EmptyPattern nil

type Reference interface{
  Next(pos Pos) (Char, Pos, bool) {
}

type Pooler interface {
  Take() (Pattern bool)
  Put(Patterns)
}

type PatternFilter func( p Pattern ) bool
type ExtenderFunc func( p Pattern, ref Reference) Patterns

func Run( ref Reference, patterns Pooler, results Pooler, 
          extender ExtenderFunc, acceptable PatternFilter) {
  patterns.Put(EmptyPattern)
  for p, valid := patterns.Take(); valid {
    pats := extender(p, ref)
    for ep := range pats {
      if acceptable( ep ) {
        patterns.Put(ep)
        results.Put(ep)
      }
    }
  }
}

func RunParallel( ref Reference, input Pooler, results Pooler, 
                  extender ExtenderFunc, acceptable PatternFilter){
  input.Put(EmptyPattern)
  start := make(chan int) // alternatively make(chan int, threadLimit)
  stop := make(chan int)
  func run {
    start <- 1
    defer func(){ stop <- 1 }()
    
    p, valid := input.Take() 
    if !valid {
      return
    }

    pats := extender(p, ref)
    for ep := range pats {
      if acceptable( ep ) {
        input.Put(ep)
        go run()
        results.Put(ep)
      }
    }
  }
  go run()
  for <-start {
    <-stop
  }
}
