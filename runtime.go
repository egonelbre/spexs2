package main

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/egonelbre/spexs2/search"
)

const mb = 1024 * 1024

func setupRuntime() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if runtime.NumCPU() < *procs || *procs <= 0 {
		*procs = runtime.NumCPU()
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
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal(err)
	}
	return true
}

func stopProfiler() {
	pprof.StopCPUProfile()
}

func setMemLimit(setup *AppSetup) {
	memLimit := uint64(*memoryLimit)
	ext := setup.Extender
	m := new(runtime.MemStats)

	setup.Extender = func(q *search.Query) search.Querys {
		runtime.ReadMemStats(m)
		if m.Alloc/mb > memLimit {
			panic("memory limit exceeded!")
		}
		return ext(q)
	}
}

func attachMemProfiler(setup *AppSetup) {
	filename := *memprofile
	limit := *memsteps

	count := 0

	ext := setup.Extender

	setup.Extender = func(q *search.Query) search.Querys {
		if count >= limit {
			f, err := os.Create(filename)
			if err != nil {
				log.Fatal(err)
			}
			_ = pprof.WriteHeapProfile(f)
			_ = f.Close()
			log.Fatal("Wrote memory profile!")
		}
		count++
		return ext(q)
	}
}

var (
	quitStats     = make(chan int)
	statsStarted  = false
	maxMemoryUsed uint64
)

func runStats(setup *AppSetup) {
	var counter uint64
	var seq string
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
			lg.Printf("%v\t%v\t%v\t%v\t%v\t\n", m.Alloc/mb, m.TotalAlloc/mb, counter, setup.Out.Len(), seq)
		}
	}()

	go func() {
		m := new(runtime.MemStats)
		t := time.Tick(50 * time.Millisecond)
		for range t {
			runtime.ReadMemStats(m)
			if m.Alloc > maxMemoryUsed {
				maxMemoryUsed = m.Alloc
			}
		}
	}()

	ext := setup.Extender
	setup.Extender = func(q *search.Query) search.Querys {
		seq = q.String()
		counter++
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
	setup.Outputtable = func(q *search.Query) bool {
		result := out(q)
		if result {
			setup.printQuery(os.Stderr, q)
		}
		return result
	}
}
