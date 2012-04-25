package main

import (
	"flag"
	"fmt"
	"runtime"
	. "spexs/trie"

	"log"
	"os"
	"runtime/pprof"
)

var (
	configs          *string = flag.String("conf", "spxs.json", "configuration file(s), comma-delimited")
	//printConf        *bool   = flag.Bool("printConf", "print the configuration file")
	details          *bool   = flag.Bool("details", false, "detailed help")
	interactiveDebug *bool   = flag.Bool("debug", false, "attach step-by-step debugger")
	verbose          *bool   = flag.Bool("verbose", false, "print extended debug info")
	version          *bool   = flag.Bool("version", false, "print version")

	procs      *int    = flag.Int("procs", 4, "processors to use")
	cpuprofile *string = flag.String("cpuprofile", "", "write cpu profile to file")
)

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
		PrintVersion()
		os.Exit(0)
	}

	initSetup()
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
