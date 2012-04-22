package debugger

import (
	"os"
	"fmt"
)

type Action interface {
	Exec(d *Debugger)
}

type Nop struct{}
type Continue struct{}
type Watch struct{}
type Break struct{}
type Disable struct{}
type Quit struct {}

type Skip struct{
	Count int
}
type Err struct {
	Msg string
}

func (n Nop) Exec(d *Debugger) {}

func (c Continue) Exec(d *Debugger) {
	d.control <- cBreak
}

func (s Skip) Exec(d *Debugger){
	d.control <- cBreak
	count := s.Count
	if count > 100 { 
		count = 100 
	}
	for i := 0; i < count; i += 1 {
		d.control <- cBreak
	}
}

func (s Watch) Exec(d *Debugger){
	d.Watch <- 1
	d.control <- cBreak
}

func (s Disable) Exec(d *Debugger){
	d.Enabled = false
	d.control <- cBreak
}

func (q Quit) Exec(d *Debugger){
	os.Exit(130)
}

func (e Err) Exec(d *Debugger){
	fmt.Fprintf(d.Logout, "err: %v\n", e.Msg)
}

func (b Break) Exec(d *Debugger){
	select {
		case <-d.Watch:
		default:
	}
}