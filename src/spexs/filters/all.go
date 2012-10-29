package filters

import . "spexs"

type CreateFunc func(Setup, []byte) Filter

var All = [...]CreateFunc{
	NoStartingGroup,
	NoEndingGroup,
}
