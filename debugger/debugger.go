package debugger

import (
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	cBreak = iota
)

type Handler func(d *Debugger, cmd string, params []string) Action

type Debugger struct {
	Commands chan string

	Enabled bool
	Watch   chan int
	Timeout int
	Handler Handler
	Logout  io.Writer

	control chan int
	lock    chan int
}

func New() *Debugger {
	d := &Debugger{}
	d.lock = make(chan int, 1)
	d.lock <- 1

	d.control = make(chan int, 10000)
	d.Commands = make(chan string, 50)
	d.Watch = make(chan int, 1)
	d.Enabled = true
	d.Timeout = 0
	d.Handler = DefaultHandler
	d.Logout = os.Stderr

	return d
}

func (d *Debugger) Break(output func()) {
	if !d.Enabled {
		return
	}

	<-d.lock

	if d.Enabled {
		output()
		d.HandleCommands()
	}

	d.lock <- 1
}

func (d *Debugger) HandleCommands() {
handling:
	for {
		select {
		case ctrl := <-d.control:
			switch ctrl {
			case cBreak:
				break handling
			}
		case cmd := <-d.Commands:
			tokens := strings.Split(cmd, " ")
			action := d.Handler(d, tokens[0], tokens[1:])
			action.Exec(d)
		case <-d.Watch:
			d.Watch <- 1
			break handling
		}
	}
}

const defaultHandlerHelp = `
  enter, c, continue: proceed to next step
  s, skip [amount]: skip breaking for [amount] (max 100)
  w, watch: watch output function without breaking
  			(to reenable breaking press enter)
  
  d, disable: disable debugger
  q, quit: terminate program
  h, help: this help
`

func DefaultHandler(d *Debugger, cmd string, params []string) Action {
	switch cmd {
	case "disable", "d":
		return Disable{}
	case "skip", "s":
		if len(params) > 0 {
			timeout, err := strconv.Atoi(params[0])
			if err == nil {
				return Skip{timeout}
			}
			return Err{err.Error()}
		}
		return Err{"Requires skip count parameter."}
	case "watch", "w":
		return Watch{}
	case "continue", "c", "":
		return Continue{}
	case "quit", "q":
		return Quit{}
	case "help", "h":
		return Msg{defaultHandlerHelp}
	}
	return Nop{}
}
