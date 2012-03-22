package main

import (
	"fmt"
	. "spexs" 
	"runtime"
	"flag"
)

var (
	//alphabet *string = flag.String("alphabet", "ACGT", "alphabet used for the input file")
    characterFile *string = flag.String("chars", "", "character set file")
    referenceFile *string = flag.String("ref", "", "reference file")   
    extenderName *string = flag.String("extender", "regex", "method used for extending nodes (simple, group, star, regex)")
    limiterName *string = flag.String("limiter", "count", "method used to determine whether node is accptable for extending (count, length, complexity)")
    limitValue *int = flag.Int("limit", 5, "value for limiter")
    topCount *int = flag.Int("top", 10, "only print top amount")
    procs *int = flag.Int("procs", 2, "processors to use")
)

var extenders = map[string] ExtenderFunc {
	"simple" : SimpleExtender,
	"group"  : GroupExtender,
	"star"   : StarExtender,
	"regex"  : GroupStarExtender,
}

type PatternFilterCreator func(limit int) PatternFilter

var limiters = map[string] PatternFilterCreator {
	"count"  : func(limit int) PatternFilter {
		return func(p Pattern) bool {
			return p.(TrieNode).Pos.Length() > limit
		}},
	"length" : func(limit int) PatternFilter {
		return func(p Pattern) bool { 
			return p.(TrieNode).Length() < limit
		}},
	"complexity" : func(limit int) PatternFilter {
		return func(p Pattern) bool { 
			return p.(TrieNode).Complexity() < limit
		}},
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
		acceptable PatternFilter
		extender ExtenderFunc
	)

	if ref, err = NewReferenceFromFile(*referenceFile, *characterFile); err != nil {
		fmt.Printf("Error occured while reading reference/character file: %v\n", err)
		return
	}
	
	extender = extenders[*extenderName]
	acceptable = limiters[*limiterName](*limitValue)

	in = NewFifoPool()
 	in.Put(*NewFullNodeFromRef(*ref))
	//out = NewFifoPool()

	f := func(a Pattern) float32{
		return float32(a.(TrieNode).Length()*a.(TrieNode).Pos.Length())
	}
	out = NewPriorityPool(f, *topCount)

	RunParallel(*ref,in,out,extender,acceptable,(*procs)*4)
	 
	p, ok := out.Take()
	for ok {
		n := p.(TrieNode)
		fmt.Printf("%s\t%v\t%v\n", n.String(), n.Pos.Length(), f(p))
		p, ok = out.Take()
	}
	fmt.Printf("\n\n")
}