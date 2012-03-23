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
    verbose *bool = flag.Bool("verbose", false, "print debug information")
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
	"count" : func(a Pattern) float32 {
		return float32(a.(*TrieNode).Pos.Len())
		},
	"complexity" : func(a Pattern) float32 {
		return float32(a.(*TrieNode).Complexity())
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

	if _, exists := fitnesses[*fitnessName]; !exists {
		fmt.Printf("Fitness function '%v' not found!\n", *limiterName)
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
	out = NewPriorityPool(fitness, *topCount)


	maxInQueue := 0
	if *verbose {
		go func(){
			for {
				fmt.Printf("queue size: %v\n", in.Len())
				if in.Len() > maxInQueue {
					maxInQueue = in.Len()
				}
				time.Sleep(1000*1000*100)
			}
		}()
	}

	if *procs == 1 {
		Run(ref,in,out,extender,acceptable)
	} else {
		RunParallel(ref,in,out,extender,acceptable,*procs)
	}	
	
	fmt.Printf("match, regexp, count, fitness\n")
	p, ok := out.Take()
	for ok {
		n := p.(*TrieNode)
		name := n.String()
		regex := ref.ReplaceGroups(name)
		fmt.Printf("%s, %v, %v, %v\n", name, regex,n.Pos.Len(), fitness(p))
		
		if *verbose {
			indices, poss := n.Pos.Iter()
			for idx := range indices {
				<-poss
				fmt.Printf("%v, ", idx)
			}
			fmt.Printf("\n\n\n")
		}
		
		p, ok = out.Take()
	}

	if *verbose {
		fmt.Printf("maximum items in queue: %v\n", maxInQueue)
	}
	
	fmt.Printf("\n")
}