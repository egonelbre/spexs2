package main

import (
	"errors"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	. "github.com/egonelbre/spexs2/search"
)

const mb = 1024 * 1024

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

type memLimiter struct {
	extender Extender
	memLimit uint64
	stats    *runtime.MemStats
}

func (e *memLimiter) Extend(q *Query) Querys {
	result := e.extender.Extend(q)
	runtime.ReadMemStats(e.stats)
	if e.stats.Alloc/mb > e.memLimit {
		panic(errors.New("MEMORY LIMIT EXCEEDED!"))
	}
	return result
}

func setMemLimit(setup *AppSetup) {
	setup.Extender = &memLimiter{setup.Extender, uint64(*memoryLimit), new(runtime.MemStats)}
}

type memProfiler struct {
	filename string
	limit    int
	count    int
	extender Extender
}

func (e *memProfiler) Extend(q *Query) Querys {
	if e.count >= e.limit {
		f, err := os.Create(e.filename)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		log.Fatal("Wrote memory profile!")
	}
	e.count += 1
	return e.extender.Extend(q)
}

func attachMemProfiler(setup *AppSetup) {
	setup.Extender = &memProfiler{*memprofile, *memsteps, 0, setup.Extender}
}

var (
	quitStats    = make(chan int)
	statsStarted = false
)

type statsCounter struct {
	db          *Database
	extender    Extender
	count       uint64
	lastPattern string
}

func (e *statsCounter) Extend(q *Query) Querys {
	e.count += 1
	e.lastPattern = q.String(e.db)
	return e.extender.Extend(q)
}

func runStats(setup *AppSetup) {
	counter := &statsCounter{setup.Db, setup.Extender, 0, ""}
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
			lg.Printf("%v\t%v\t%v\t%v\t%v\t\n", m.Alloc/mb, m.TotalAlloc/mb, counter.count, setup.Out.Len(), counter.lastPattern)
		}
	}()

	setup.Extender = counter
}

func endStats() {
	if statsStarted {
		quitStats <- 1
	}
}

type liveViewer struct {
	setup  *AppSetup
	filter Filter
}

func (f *liveViewer) Accepts(q *Query) bool {
	ok := f.filter.Accepts(q)
	if ok {
		f.setup.printQuery(os.Stderr, q)
	}
	return ok
}

func setupLiveView(setup *AppSetup) {
	setup.Outputtable = &liveViewer{setup, setup.Outputtable}

}
