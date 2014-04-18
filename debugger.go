package main

import (
	"fmt"

	"github.com/egonelbre/spexs2/debugger"
	. "github.com/egonelbre/spexs2/search"
)

var dbg = debugger.New()

type debugExtender struct {
	Db          *Database
	extender    Extender
	extendable  Filter
	outputtable Filter
}

func (e *debugExtender) Extend(q *Query) Querys {
	result := e.extender.Extend(q)
	dbg.Break(func() {
		fmt.Fprintf(dbg.Logout, "extending: %v\n", q.String(e.Db))
		for _, extended := range result {
			fmt.Fprintf(dbg.Logout, "  | %v\n", q.String(e.Db))
			if *verbose {
				fmt.Fprintf(dbg.Logout, "\tE:%v\tO:%v\n",
					e.extendable.Accepts(extended),
					e.outputtable.Accepts(extended))
			}
		}
	})
	return result
}

func attachDebugger(s *AppSetup) {
	debugger.RunShell(dbg)
	s.Extender = &debugExtender{s.Db, s.Extender, s.Extendable, s.Outputtable}
}
