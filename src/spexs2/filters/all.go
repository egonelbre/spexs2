package filters

import . "spexs"

type CreateFunc func(Setup, []byte) (Filter, error)

var All = [...]CreateFunc{
	NoGroupingEnds,
}
