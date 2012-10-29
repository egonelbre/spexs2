package filters

import . "spexs"

type CreateFunc func(Setup, []byte) Filter

var All = [...]CreateFunc{
	NoStartingGroup,
	NoEndingGroup,
}

var Help = `
:Pattern:
  NoStartingGroup : removes patterns with starting group token
  NoEndingGroup   : removes patterns with ending group token
  	                (useful only in output.filter)

:Feature:
  Any feature can be used as a filter
`
