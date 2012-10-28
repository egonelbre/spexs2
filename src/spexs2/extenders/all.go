package extenders

import . "spexs"

type CreateFunc func(Setup, []byte) Extender

var All = [...]Extender{
	Simplex,
	Groupex,
	Starex,
	Regex,
}

func wrap(f Extender) CreateFunc {
	return func(s Setup, data []byte) Extender {
		return f
	}
}
