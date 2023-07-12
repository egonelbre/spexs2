package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/egonelbre/spexs2/search"
)

var (
	printVersion = flag.Bool("version", false, "print version")

	details = flag.Bool("details", false, "detailed help")
	verbose = flag.Bool("verbose", false, "print extended info")

	interactiveDebug = flag.Bool("debug", false, "attach step-by-step debugger")
	live             = flag.Bool("live", false, "print live output")
	configs          = flag.String("conf", "", "configuration file(s), comma-delimited")
	writeConf        = flag.String("writeconf", "", "write conf file")

	stats       = flag.Bool("stats", false, "print memory/extension statistics")
	procs       = flag.Int("procs", -1, "goroutines for extending")
	memoryLimit = flag.Int("mem", -1, "memory limit in MB")

	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile = flag.String("memprofile", "", "write mem profile to file")
	memsteps   = flag.Int("memsteps", 10000, "after how many extensions to write the mem profile")

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

	if *printVersion {
		PrintVersion()
		os.Exit(0)
	}

	conf := NewConf(*configs)

	if *writeConf != "" {
		conf.WriteToFile(*writeConf)
		return
	}

	if *configs == "" {
		log.Fatal("Configuration file is required!")
		return
	}

	info("reading input")

	setup := NewAppSetup(conf)

	// defined in runtime.go
	setupRuntime()

	if startProfiler(*cpuprofile) {
		defer stopProfiler()
	}

	ifthen := func(val bool, f func(*AppSetup)) {
		if val {
			f(setup)
		}
	}

	ifthen(*interactiveDebug, attachDebugger)
	ifthen(*stats, runStats)
	ifthen(*live, setupLiveView)
	ifthen(*memoryLimit > 0, setMemLimit)
	ifthen(*memprofile != "", attachMemProfiler)

	info("running spexs [", *procs, "]")

	starttime := time.Now()

	if *procs == 1 {
		search.Run(&setup.Setup)
	} else {
		search.RunParallel(&setup.Setup, *procs)
	}

	finishtime := time.Now()

	endStats()

	info("printing results")
	setup.Printer(os.Stdout, setup.Out)

	info("searching ", finishtime.Sub(starttime))
	if *stats {
		info("max ", maxMemoryUsed/(1024*1024), "MB")
	}
}
