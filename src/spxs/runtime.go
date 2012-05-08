package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	. "spexs"
	"time"
)

const mb = 1024 * 1024

func setupRuntime() {
	runtime.GOMAXPROCS(runtime.NumCPU())
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

func setMemLimit(setup *AppSetup, memLimit uint64) {
	ext := setup.Extender
	m := new(runtime.MemStats)

	setup.Extender = func(p *Pattern, ref *Reference) Patterns {
		runtime.ReadMemStats(m)
		if m.Alloc/mb > memLimit {
			panic(errors.New("MEMORY LIMIT EXCEEDED!"))
		}
		return ext(p, ref)
	}
}

var quitStats = make(chan int)

func runStats(setup *AppSetup) {
	var counter uint64 = 0
	var seq string = ""

	go func() {
		m := new(runtime.MemStats)
		for {
			select {
			case <-time.After(200 * time.Millisecond):
			case <-quitStats:
				return
			}

			runtime.ReadMemStats(m)
			fmt.Fprintf(os.Stderr, "%v\t%v\t%v\t%v\t%s\n", runtime.NumGoroutine(), m.Alloc/mb, m.TotalAlloc/mb, counter, seq)
		}
	}()

	ext := setup.Extender
	setup.Extender = func(p *Pattern, ref *Reference) Patterns {
		seq = p.String()
		counter += 1
		return ext(p, ref)
	}
}

func endStats() {
	quitStats <- 1
}

func setupLiveView(setup *AppSetup) {
	out := setup.Outputtable
	setup.Outputtable = func(p *Pattern, ref *Reference) bool {
		result := out(p, ref)
		if result {
			setup.Printer(os.Stderr, p, ref)
		}
		return result
	}
}
