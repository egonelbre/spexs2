package filters

import (
	. "spexs"
	"utils"
)

type CreateFunc func(Setup, []byte) Filter

var All = [...]CreateFunc{
	NoStartingGroup,
	NoEndingGroup,
}

func Get(name string) (CreateFunc, bool) {
	for _, fn := range All {
		if utils.FuncName(fn) == name {
			return fn, true
		}
	}
	return nil, false
}

var Help = `
:Pattern:
  NoStartingGroup : removes patterns with starting group token
  NoEndingGroup   : removes patterns with ending group token
  	                (useful only in output.filter)

:Feature:
  Any feature can be used as a filter
`
