package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	. "spexs"
	"spexs/pool"
	"time"
)

var (
	version *bool = flag.Bool("version", false, "print version")

	details *bool = flag.Bool("details", false, "detailed help")
	verbose *bool = flag.Bool("verbose", false, "print extended info")

	interactiveDebug *bool   = flag.Bool("debug", false, "attach step-by-step debugger")
	live             *bool   = flag.Bool("live", false, "print live output")
	configs          *string = flag.String("conf", "", "configuration file(s), comma-delimited")

	stats       *bool = flag.Bool("stats", false, "print memory/extension statistics")
	procs       *int  = flag.Int("procs", 16, "goroutines for extending")
	memoryLimit *int  = flag.Int("mem", -1, "memory limit in MB")

	cpuprofile *string = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile *string = flag.String("memprofile", "", "write mem profile to file")
	memsteps   *int    = flag.Int("memsteps", 10000, "after how many extensions to write the mem profile")

	// logging to stderr
	lg = log.New(os.Stderr, "", log.Ltime)
)

func info(v ...interface{}) {
	if *verbose {
		fmt.Fprintf(os.Stderr, "%+v: ", time.Now())
		fmt.Fprintln(os.Stderr, v...)
	}
}

func main() {
	flag.Parse()

	if *details {
		PrintHelp()
		os.Exit(0)
	}

	if *version {
		PrintVersion()
		os.Exit(0)
	}

	if *configs == "" {
		fmt.Fprintf(os.Stderr, "Configuration file is required!\n")
		return
	}

	conf := ReadConfiguration(*configs)

	info("reading input")

	setup := CreateSetup(conf)

	// defined in runtime.go
	setupRuntime()

	if startProfiler(*cpuprofile) {
		defer stopProfiler()
	}

	ifthen := func(val bool, f func(*AppSetup)) {
		if val {
			f(&setup)
		}
	}

	ifthen(*interactiveDebug, attachDebugger)
	ifthen(*stats, runStats)
	ifthen(*live, setupLiveView)
	ifthen(*memoryLimit > 0, setMemLimit)
	ifthen(*memprofile != "", attachMemProfiler)

	info("running spexs")

	if *procs <= 1 {
		Run(&setup.Setup)
	} else {
		RunParallel(&setup.Setup, *procs)
	}

	endStats()

	limit := conf.Output.Count
	setup.Printer(os.Stdout, nil)

	info("getting best results")

	best := setup.Out.(*pool.Priority).Bottom(limit)

	info("printing results")
	for _, node := range best {
		setup.Printer(os.Stdout, node)
	}

	info("done")

	/*

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
		}*/
}
