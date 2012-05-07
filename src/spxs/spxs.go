package main

import (
	"flag"
	"fmt"
	"runtime"
	. "spexs"

	"log"
	"os"
	"runtime/pprof"
	"time"
	"errors"
)

var (
	configs          *string = flag.String("conf", "spxs.json", "configuration file(s), comma-delimited")
	//printConf        *bool   = flag.Bool("printConf", "print the configuration file")
	memoryLimit      *int    = flag.Int("mem", 1024, "memory limit in MB")
	details          *bool   = flag.Bool("details", false, "detailed help")
	interactiveDebug *bool   = flag.Bool("debug", false, "attach step-by-step debugger")
	verbose          *bool   = flag.Bool("verbose", false, "print extended debug info")
	version          *bool   = flag.Bool("version", false, "print version")
	live             *bool   = flag.Bool("live", false, "print output live")

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

	fmt.Fprintf(os.Stderr, "reading input\n")

	initSetup()
	setup := CreateSetup(conf)

	setupRuntime()

	if startProfiler(*cpuprofile) {
		defer stopProfiler()
	}

	if *interactiveDebug {
		AttachDebugger(&setup)
	}

	var counter uint64 = 0
	var seq string = ""

	if *verbose {
		go func(){
		m := new(runtime.MemStats)
		gb := uint64(1024*1024)
		for {
			runtime.ReadMemStats(m)
			fmt.Printf("%v\t%v\t%v\t%v\t%s\n", runtime.NumGoroutine(), m.Alloc/gb, m.TotalAlloc/gb, counter, seq)
			time.Sleep(200 * time.Millisecond)

			if m.Alloc/gb > uint64(*memoryLimit) {
				panic(errors.New("MEMORY LIMIT EXCEEDED!"))
			}
			}
		}()
		
		ext := setup.Extender
		setup.Extender = func(p *Pattern, ref *Reference) Patterns {
			seq = p.String()
			counter += 1
			return ext(p, ref)
		}
	}

	if *live {
		out := setup.Outputtable
		setup.Outputtable = func(p *Pattern, ref *Reference) bool {
			result := out(p, ref)
			if result {
				setup.Printer(os.Stderr, p, ref)
			}
			return result
		}
	}

	fmt.Fprintf(os.Stderr, "running spexs\n")

	if *procs <= 1 {
		Run(&setup.Setup)
	} else {
		RunParallel(&setup.Setup, *procs)
	}

	setup.Printer(os.Stdout, nil, setup.Ref)

	if conf.Output.Queue != "lifo"  {
		fmt.Fprintf(os.Stderr, "throwing away bad results\n")

		limit := conf.Output.Count
		if limit > 0 {
			for setup.Out.Len() > limit {
				setup.Out.Take()
			}
		}
	}
	
	fmt.Fprintf(os.Stderr, "printing results\n")

	node, ok := setup.Out.Take()
	for ok {
		setup.Printer(os.Stdout, node, setup.Ref)
		node, ok = setup.Out.Take()
	}

	fmt.Printf("\n")
}
