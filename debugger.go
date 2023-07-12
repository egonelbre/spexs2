package main

import (
	"fmt"

	"github.com/egonelbre/spexs2/debugger"
	"github.com/egonelbre/spexs2/search"
)

var dbg = debugger.New()

func attachDebugger(s *AppSetup) {
	debugger.RunShell(dbg)
	f := s.Extender
	s.Extender = func(q *search.Query) search.Querys {
		result := f(q)
		dbg.Break(func() {
			fmt.Fprintf(dbg.Logout, "extending: %v\n", q.String())
			for _, extended := range result {
				fmt.Fprintf(dbg.Logout, "  | %v\n", q.String())
				if *verbose {
					fmt.Fprintf(dbg.Logout, "\tE:%v\tO:%v\n",
						s.Extendable(extended),
						s.Outputtable(extended))
				}
			}
		})
		return result
	}
}
