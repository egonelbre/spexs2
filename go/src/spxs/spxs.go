package main

import (
	"fmt"
	. "spexs" 
	"runtime"
	"time"
	"flag"

	"runtime/pprof"
	"os"
	"log"
)

/*
	multiple limiters
	output limiters
	p-value binom/hyper

	flexibility wildcards
	[-max_gap_nr nr]		- How many flexible gaps at most
 	[-min_gap nr] 			- minimum length of a gap
 	[-max_gap nr] 			- maximum length of a gap
 	[-no_gap_len nr] 		- require at least so many positions gap-less
 	[-init_gap_len nr] 		- Initiate that value (can require longer/shorter first motif...)
 	[-only_print_if_gap_allowed]	- only report motif if gap could be allowed at that pos
	
	output ===
	-length:6..
	-count:10..

	fitness ===
	-p-value: -1


	/group/inbox/elbre

 	// -acceptable
*/

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
    cpuprofile *string = flag.String("cpuprofile", "", "write cpu profile to file")
)

var extenders = map[string] TrieExtenderFunc {
	"simple" : SimpleExtender,
	"group"  : GroupExtender,
	"star"   : StarExtender,
	"regexp"  : GroupStarExtender,
}

type PatternFilterCreator func(limit int) TrieFilterFunc

var limiters = map[string] PatternFilterCreator {
	"count"  : func(limit int) TrieFilterFunc {
		return func(p *TrieNode) bool {
			return p.Pos.Len() >= limit
		}},
	"length" : func(limit int) TrieFilterFunc {
		return func(p *TrieNode) bool { 
			return p.Len() <= limit
		}},
	"complexity" : func(limit int) TrieFilterFunc {
		return func(p *TrieNode) bool { 
			return p.Complexity() <= limit
		}},
}

var fitnesses = map[string] TrieFitnessFunc {
	"def" : func(p *TrieNode) float32 {
		return float32(p.Len()*p.Pos.Len())
		},
	"len" : func(p *TrieNode) float32 {
		return float32(p.Len())
		},
	"count" : func(p *TrieNode) float32 {
		return float32(p.Pos.Len())
		},
	"complexity" : func(p *TrieNode) float32 {
		return float32(p.Complexity())
		},
}

func inputOrdering(p *TrieNode) float32 {
	return 1/float32(p.Len())
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

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var (
		ref *UnicodeReference
		out TriePooler
		in TriePooler
		acceptable TrieFilterFunc
		extender TrieExtenderFunc
		fitness TrieFitnessFunc
	)

	ref = NewReferenceFromFile(*referenceFile, *characterFile)
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
		RunTrie(ref,in,out,extender,acceptable)
	} else {
		RunTrieParallel(ref,in,out,extender,acceptable,*procs)
	}	
	
	fmt.Printf("match, regexp, count, fitness\n")
	node, ok := out.Take()
	for ok {
		name := node.String()
		regex := ref.ReplaceGroups(name)
		fmt.Printf("%s, %v, %v, %v\n", name, regex, node.Pos.Len(), fitness(node))
		
		if *verbose {
			for idx := range node.Pos.Iter() {
				fmt.Printf("%v, ", idx)
			}
			fmt.Printf("\n\n\n")
		}
		
		node, ok = out.Take()
	}

	if *verbose {
		fmt.Printf("maximum items in queue: %v\n", maxInQueue)
	}
	
	fmt.Printf("\n")
}