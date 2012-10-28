package filters

type CreateFunc func(Conf, Setup) (Func, error)

var All = [...]CreateFunc{
	NoGroupingEnds,
}
