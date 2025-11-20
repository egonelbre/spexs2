package debugger

import (
	"fmt"
	"os"
)

type Action interface {
	Exec(d *Debugger)
}

type Nop struct{}
type Continue struct{}
type Watch struct{}
type Disable struct{}
type Quit struct{}

type Skip struct {
	Count int
}
type Msg struct {
	Msg string
}
type Err struct {
	Msg string
}

func (n Nop) Exec(d *Debugger) {}

func (c Continue) Exec(d *Debugger) {
	d.control <- cBreak
	select {
	case <-d.Watch:
	default:
	}
}

func (s Skip) Exec(d *Debugger) {
	d.control <- cBreak
	count := min(s.Count, 100)
	for i := 0; i < count; i++ {
		d.control <- cBreak
	}
}

func (s Watch) Exec(d *Debugger) {
	d.Watch <- 1
	d.control <- cBreak
}

func (s Disable) Exec(d *Debugger) {
	d.Enabled = false
	d.control <- cBreak
}

func (q Quit) Exec(d *Debugger) {
	os.Exit(130)
}

func (e Err) Exec(d *Debugger) {
	fmt.Fprintf(d.Logout, "err: %v\n", e.Msg)
}

func (m Msg) Exec(d *Debugger) {
	fmt.Fprintf(d.Logout, "%v\n", m.Msg)
}
