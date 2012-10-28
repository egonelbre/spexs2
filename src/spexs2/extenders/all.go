package extenders

import . "spexs"

type CreateFunc func(Conf, Setup) (Func, error)

var All = [...]CreateFunc{
	wrap(Simplex),
	wrap(Groupex),
	wrap(Starex),
	wrap(Regex),
}

func wrap(f Extender) CreateFunc {
	return func(conf Conf) (Func, error) {
		return f, nil
	}
}
