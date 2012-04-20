package debugger

type Handler func(d *Debugger, cmd string, params []string)

func DefaultHandler(d *Debugger, cmd string, params []string ){
	switch cmd {
		"info", "i":;
		"run", "r":;
		"skip", "s":;
		"continue", "c", "":;
		"quit", "q" :;
	}
}

type Debugger struct{
	commands chan string
	Enabled bool
	Timeout int
	Handler Handler

	lock chan int
}

func New(f Handler) *Debugger {
	d := &Debugger{}
	d.lock = make(chan int, 1)
	d.lock <- 1

	d.Commands = make(chan string, 50)
	d.Enabled = true
	d.Timeout = 0
	d.Handler = f
}

func (d *Debugger) Debug(){
	d.DebugPrint(func(){})
}

func (d *Debugger) DebugPrint( output func() ){
	if !d.Enabled {
		return
	}

	d.Lock()

	output()

	d.HandleCommands()
	d.Unlock()
}

func (d *Debugger) HandleCommands(){
	cmd <- d.commands
	d.Handler(cmd, [])
}

func (d *Debugger) Lock(){
	<- d.lock
}

func (d *Debugger) Unlock(){
	d.lock <- 1
}
