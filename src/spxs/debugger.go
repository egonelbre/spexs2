package main

import (
	"debugger"
	"fmt"
	. "spexs"
)

var dbg = debugger.New()

func attachDebugger(s *AppSetup) {
	debugger.RunShell(dbg)
	f := s.Extender
	s.Extender = func(q *Query, db *Database) Patterns {
		tmp := f(q, db)
		result := NewPatterns()
		dbg.Break(func() {
			fmt.Fprintf(dbg.Logout, "extending: %v\n", q.ToString(db))
			for extended := range tmp {
				result <- extended
				fmt.Fprintf(dbg.Logout, "  | %v\n", q.ToString(db))
				if *verbose {
					fmt.Fprintf(dbg.Logout, "      E:%v  O:%v\n", s.Extendable(extended, ref), s.Outputtable(extended, ref))
					fmt.Fprintf(dbg.Logout, "      ")
					s.Printer(dbg.Logout, extended, db)
				}
			}
			close(result)
		})
		return result
	}
}
