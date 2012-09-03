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
	s.Extender = func(p *Query, ref *Database) Patterns {
		tmp := f(p, ref)
		result := NewPatterns()
		dbg.Break(func() {
			fmt.Fprintf(dbg.Logout, "extending: %v\n", p.String())
			for extended := range tmp {
				result <- extended
				fmt.Fprintf(dbg.Logout, "  | %v\n", ref.ReplaceGroups(extended.String()))
				if *verbose {
					fmt.Fprintf(dbg.Logout, "      E:%v  O:%v\n", s.Extendable(extended, ref), s.Outputtable(extended, ref))
					fmt.Fprintf(dbg.Logout, "      ")
					s.Printer(dbg.Logout, extended, ref)
				}
			}
			close(result)
		})
		return result
	}
}
