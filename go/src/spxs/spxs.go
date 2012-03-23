package main

import (
	"fmt"
	. "spexs" 
	"runtime"
	"time"
	"flag"
)

var (
    characterFile *string = flag.String("chars", "", "character set file")
    referenceFile *string = flag.String("ref", "", "reference file")   
    extenderName *string = flag.String("extender", "regexp", "method used for extending nodes (simple, group, star, regex)")
    limiterName *string = flag.String("limiter", "count", "method used to determine whether node is accptable for extending (count, length, complexity)")
    fitnessName *string = flag.String("fitness", "def", "fitness function used for sorting (def)")
    limitValue *int = flag.Int("limit", 5, "value for limiter")
    topCount *int = flag.Int("top", 10, "only print top amount")
    procs *int = flag.Int("procs", 4, "processors to use")
)

var extenders = map[string] ExtenderFunc {
	"simple" : SimpleExtender,
	"group"  : GroupExtender,
	"star"   : StarExtender,
	"regexp"  : GroupStarExtender,
}

type PatternFilterCreator func(limit int) FilterFunc

var limiters = map[string] PatternFilterCreator {
	"count"  : func(limit int) FilterFunc {
		return func(p Pattern) bool {
			return p.(*TrieNode).Pos.Len() >= limit
		}},
	"length" : func(limit int) FilterFunc {
		return func(p Pattern) bool { 
			return p.(*TrieNode).Len() <= limit
		}},
	"complexity" : func(limit int) FilterFunc {
		return func(p Pattern) bool { 
			return p.(*TrieNode).Complexity() <= limit
		}},
}

var fitnesses = map[string] FitnessFunc {
	"def" : func(a Pattern) float32 {
		return float32(a.(*TrieNode).Len()*a.(*TrieNode).Pos.Len())
		},
	"len" : func(a Pattern) float32 {
		return float32(a.(*TrieNode).Len())
		},
}

func inputOrdering(a Pattern) float32 {
	return 1/float32(a.(*TrieNode).Len())
}

func main() {
	flag.Parse()

	ok := true

	if *referenceFile == "" || *characterFile == "" {
		fmt.Printf("Reference and character files are required!\n")
		ok = false
	}

	if _, exists := extenders[*extenderName]; !exists {
		fmt.Printf("Extender function '%v' not found!\n", *extenderName)
		ok = false
	}

	if _, exists := limiters[*limiterName]; !exists {
		fmt.Printf("Limiter function '%v' not found!\n", *limiterName)
		ok = false
	}

	if !ok {
		return
	}
	
	if *procs > 0 {
		runtime.GOMAXPROCS(*procs)
	}

	var (
		ref *UnicodeReference
		err error
		out Pooler
		in Pooler
		acceptable FilterFunc
		extender ExtenderFunc
		fitness FitnessFunc
	)

	if ref, err = NewReferenceFromFile(*referenceFile, *characterFile); err != nil {
		fmt.Printf("Error occured while reading reference/character file: %v\n", err)
		return
	}
	
	extender = extenders[*extenderName]
	acceptable = limiters[*limiterName](*limitValue)
	fitness = fitnesses[*fitnessName]

	in = NewPriorityPool(inputOrdering, 1000000000)
 	in.Put(NewFullNodeFromRef(ref))
 	out = NewFifoPool()
	//out = NewPriorityPool(fitness, *topCount)


	maxInQueue := 0
	go func(){
		for {
			print(in.Len(), "\n")
			if in.Len() > maxInQueue {
				maxInQueue = in.Len()
			}
			time.Sleep(1000*1000*100)
		}
	}()

	RunParallel(ref,in,out,extender,acceptable,(*procs)*4)
	 
	p, ok := out.Take()
	for ok {
		n := p.(*TrieNode)
		name := ref.ReplaceGroups(n.String())
		fmt.Printf("%s\t%v\t%v\n", name, n.Pos.Len(), fitness(p))
		p, ok = out.Take()
	}

	print("maximum items in queue: ", maxInQueue, "\n")
	fmt.Printf("\n\n")
}