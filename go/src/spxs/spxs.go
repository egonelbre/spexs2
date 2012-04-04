package main

import (
	"flag"
	"fmt"
	"runtime"
	. "spexs"
	"time"

	"log"
	"os"
	"runtime/pprof"
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
	configs *string = flag.String("conf", "spxs.json", "configuration file(s), comma-delimited")

	procs      *int    = flag.Int("procs", 4, "processors to use")
	verbose    *bool   = flag.Bool("verbose", false, "print debug information")
	cpuprofile *string = flag.String("cpuprofile", "", "write cpu profile to file")
)

func setupRuntime() {
	
}

func main() {
	flag.Parse()

	if *configs == "" {
		fmt.Println("Configuration file is required!")
		return
	}

	conf := ReadConfiguration(*configs)
	setup := CreateSetup(conf)

	// RUNTIME SETUP

	if *procs > 0 {
		runtime.GOMAXPROCS(*procs)
	}

	// DEBUG SETUP

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// debugging the input queue
	if *verbose {
		go func() {
			for {
				fmt.Printf("queue size: %v\n", setup.In.Len())
				time.Sleep(1000 * 1000 * 100)
			}
		}()
	}	

	if *procs == 1 {
		RunTrie(setup.Ref, setup.In, setup.Out, setup.Extender, setup.Extendable, setup.Outputtable)
	} else {
		RunTrieParallel(setup.Ref, setup.In, setup.Out, setup.Extender, setup.Extendable, setup.Outputtable, *procs)
	}

	fmt.Printf("match, regexp, count, fitness, p-value\n")
	node, ok := setup.Out.Take()
	for ok {
		name := node.String()
		regex := setup.Ref.ReplaceGroups(name)
		fmt.Printf("%s, %v, %v, %v, %v\n", name, regex, node.Pos.Len(), setup.Fitness(node), node.PValue(setup.Ref))

		if *verbose {
			for idx := range node.Pos.Iter() {
				fmt.Printf("%v, ", idx)
			}
			fmt.Printf("\n\n\n")
		}

		node, ok = setup.Out.Take()
	}

	fmt.Printf("\n")
}
