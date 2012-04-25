package main

import (
	"flag"
	"fmt"
	"runtime"
	. "spexs/trie"

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
	SPXS_VERSION = "0.1"
	configs *string = flag.String("conf", "spxs.json", "configuration file(s), comma-delimited")
	details *bool   = flag.Bool("details", false, "detailed help")
	interactiveDebug *bool  = flag.Bool("debug", false, "attach step-by-step debugger")
	verbose    *bool   = flag.Bool("verbose", false, "print extended debug info")
	version    *bool   = flag.Bool("version", false, "print version")

	procs      *int    = flag.Int("procs", 4, "processors to use")
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

func PrintVersion(conf Conf) {
	fmt.Fprintf(os.Stderr, "%v\n", SPXS_VERSION)
}

func setupRuntime() {
	if *procs > 0 {
		runtime.GOMAXPROCS(*procs)
	}
}

func startProfiler(outputFile string) bool {
	if outputFile == "" {
		return false
	}
	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	return true
}

func stopProfiler() {
	pprof.StopCPUProfile()
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

	if *version {
		PrintVersion(conf)
		os.Exit(0)
	}

	setup := CreateSetup(conf)

	setupRuntime()

	if startProfiler(*cpuprofile) {
		defer stopProfiler()
	}

	if *interactiveDebug {
		AttachDebugger(&setup)
	}

	if *procs <= 1 {
		Run(setup.Setup)
	} else {
		RunParallel(setup.Setup, *procs)
	}

	setup.Printer(os.Stdout, nil, setup.Ref)

	node, ok := setup.Out.Take()
	for ok {
		setup.Printer(os.Stdout, node, setup.Ref)
		node, ok = setup.Out.Take()
	}

	fmt.Printf("\n")
}
