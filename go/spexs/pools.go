package spexs

type FifoPool struct {
  patterns Patterns
}

func (p *FifoPool) Take() (Pattern, bool) {
  return <- p.patterns
}

func (p *FifoPool) Put( pat Pattern ) {
  p.patterns <- pat
}

import "container/heap"

type FitnessFunction func(a Pattern) float

type PriorityPool struct {
	token chan int
	input Patterns
	items []Pattern
	Fitness FitnessFunction
}

func NewPriorityPool() *PriorityPool {
	p := &PriorityPool{}
	p.token = make(chan int)
	p.input = make(Patterns)
	p.items = make([]*Pattern)
	token <- 1

	go func(){ 
		// put synchronizer
		for c := <- p.input {
			p.actualPut(c)
		}
	}
}

func (p *PriorityPool) Take() (Pattern, bool) {
	<- token
	defer func(){ token <- 1 }

  	v := heap.Pop(p, pat)
  	return v.(Pattern)
}

func (p *PriorityPool) actualPut( pat Pattern ) {
	<- token
	defer func(){ token <- 1 }

	heap.Push(p, pat)
}

func (p *PriorityPool) Put( pat Pattern ) {
	p.input <- pat
}

// sort.Interface
func (p *PriorityPool) Len(){
	return len(p.items)
}

func (p *PriorityPool) Swap(i, j int){
	temp := p.items[i]
	p.items[i] = p.items[j]
	p.items[j] = temp
}

func (p *PriorityPool) Less(i, j int) bool{
	return p.Fitness(p.items[i]) < p.Fitness(p.items[j])
}

// heap.Interface
func (p *PriorityPool) Push(x interface{}) {
	append( p.items, x.(Pattern) )
}

func (p *PriorityPool) Pop() interface{} {
	r := p.items[-1]
	p.items = p.items[0:len(p.items)-1]
	append( p.items, x.(Pattern) )
}
