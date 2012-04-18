package main

import (
	"flag"
	"fmt"
	"runtime"
	. "spexs/trie"
	"time"

	"log"
	"os"
	"runtime/pprof"
	"sort"
)

/*
	p-value binom

	flexibility wildcards
	[-max_gap_nr nr]		- How many flexible gaps at most
 	[-min_gap nr] 			- minimum length of a gap
 	[-max_gap nr] 			- maximum length of a gap
 	[-no_gap_len nr] 		- require at least so many positions gap-less
 	[-init_gap_len nr] 		- Initiate that value (can require longer/shorter first motif...)
 	[-only_print_if_gap_allowed]	- only report motif if gap could be allowed at that pos

*/

var (
	configs *string = flag.String("conf", "spxs.json", "configuration file(s), comma-delimited")
	details *bool   = flag.Bool("details", false, "detailed help")

	procs      *int    = flag.Int("procs", 4, "processors to use")
	verbose    *bool   = flag.Bool("verbose", false, "print debug information")
	cpuprofile *string = flag.String("cpuprofile", "", "write cpu profile to file")
)

func PrintHelp(conf Conf) {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "spxs [FLAGS] [OPTIONS]\n\n")
	fmt.Fprintf(os.Stderr, "FLAGS: \n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nOPTIONS: \n")

	keys := make([]string, len(conf.Aliases))
	i := 0
	for key := range conf.Aliases {
		keys[i] = key
		i += 1
	}
	sort.Strings(keys)

	for _, name := range keys {
		args := conf.Aliases[name]
		fmt.Fprintf(os.Stderr, "  %s : %s\n", name, args.Desc)
	}

	fmt.Fprintf(os.Stderr, "\nEXAMPLES: \n")
	fmt.Fprintf(os.Stderr, "  spxs -procs=4 ref=data.dna val=random.dna\n")

	fmt.Fprintf(os.Stderr, "\n")
}

func main() {
	flag.Parse()

	if *configs == "" {
		fmt.Println("Configuration file is required!")
		return
	}

	conf := ReadConfiguration(*configs)

	if *details {
		PrintHelp(conf)
		os.Exit(0)
	}

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

	if *procs <= 1 {
		Run(setup.Setup)
	} else {
		RunParallel(setup.Setup, *procs)
	}


	node, ok := setup.Out.Take()
	for ok {
		setup.Printer(os.Stdout, node, setup.Ref)
		node, ok = setup.Out.Take()
	}
	/*fmt.Printf("match, regexp, count, fitness, p-value\n")
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
	}*/

	fmt.Printf("\n")
}
