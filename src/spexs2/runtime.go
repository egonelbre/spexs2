package main

import (
	"errors"
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

func setMemLimit(setup *AppSetup) {
	memLimit := uint64(*memoryLimit)
	ext := setup.Extender
	m := new(runtime.MemStats)

	setup.Extender = func(q *Query) Querys {
		runtime.ReadMemStats(m)
		if m.Alloc/mb > memLimit {
			panic(errors.New("MEMORY LIMIT EXCEEDED!"))
		}
		return ext(q)
	}
}

func attachMemProfiler(setup *AppSetup) {
	filename := *memprofile
	profile := 0
	count := 0
	limit := *memsteps

	ext := setup.Extender

	setup.Extender = func(q *Query) Querys {
		if count == limit {
			f, err := os.Create(filename + string(profile))
			if err != nil {
				log.Fatal(err)
			}
			pprof.WriteHeapProfile(f)
			f.Close()
			lg.Printf("Wrote memory profile %v!\n", profile)

			profile += 1
			count = 0
		}
		count += 1
		return ext(q)
	}
}

var (
	quitStats    = make(chan int)
	statsStarted = false
)

func runStats(setup *AppSetup) {
	var counter uint64 = 0
	var seq string = ""
	statsStarted = true

	go func() {
		m := new(runtime.MemStats)
		for {
			select {
			case <-time.After(200 * time.Millisecond):
			case <-quitStats:
				return
			}

			runtime.ReadMemStats(m)
			lg.Printf("%v\t%v\t%v\t%v\n", m.Alloc/mb, m.TotalAlloc/mb, counter, seq)
		}
	}()

	ext := setup.Extender
	setup.Extender = func(q *Query) Querys {
		seq = q.String()
		counter += 1
		return ext(q)
	}
}

func endStats() {
	if statsStarted {
		quitStats <- 1
	}
}

func setupLiveView(setup *AppSetup) {
	out := setup.Outputtable
	setup.Outputtable = func(q *Query) bool {
		result := out(q)
		if result {
			setup.Printer(os.Stderr, q)
		}
		return result
	}
}
