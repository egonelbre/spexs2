package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	. "spexs"
)

var (
	version *bool = flag.Bool("version", false, "print version")

	details *bool = flag.Bool("details", false, "detailed help")
	verbose *bool = flag.Bool("verbose", false, "print extended info")

	interactiveDebug *bool   = flag.Bool("debug", false, "attach step-by-step debugger")
	live             *bool   = flag.Bool("live", false, "print live output")
	configs          *string = flag.String("conf", "spxs.json", "configuration file(s), comma-delimited")

	stats       *bool = flag.Bool("stats", false, "print memory/extension statistics")
	procs       *int  = flag.Int("procs", 4, "parallel routines for extending")
	memoryLimit *int  = flag.Int("mem", -1, "memory limit in MB")

	cpuprofile *string = flag.String("cpuprofile", "", "write cpu profile to file")

	// logging to stderr
	lg = log.New(os.Stderr, "", log.Ltime)
)

func info(v ...interface{}) {
	if *verbose {
		fmt.Fprintln(os.Stderr, v...)
	}
}

func main() {
	flag.Parse()

	if *configs == "" {
		fmt.Fprintf(os.Stderr, "Configuration file is required!")
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

	info("reading input")

	initSetup()
	setup := CreateSetup(conf)

	// defined in runtime.go
	setupRuntime()

	if startProfiler(*cpuprofile) {
		defer stopProfiler()
	}

	if *interactiveDebug {
		attachDebugger(&setup)
	}

	if *stats {
		runStats(&setup)
	}

	if *live {
		setupLiveView(&setup)
	}

	if *memoryLimit > 0 {
		setMemLimit(&setup, uint64(*memoryLimit))
	}

	info("running spexs")

	if *procs <= 1 {
		Run(&setup.Setup)
	} else {
		RunParallel(&setup.Setup, *procs)
	}

	endStats()

	if conf.Output.Queue != "lifo" {
		info("sorting results")

		limit := conf.Output.Count
		if limit > 0 {
			for setup.Out.Len() > limit {
				setup.Out.Take()
			}
		}
	}

	info("printing results")

	setup.Printer(os.Stdout, nil, setup.Ref)
	node, ok := setup.Out.Take()
	for ok {
		setup.Printer(os.Stdout, node, setup.Ref)
		node, ok = setup.Out.Take()
	}

	fmt.Printf("\n")
}
