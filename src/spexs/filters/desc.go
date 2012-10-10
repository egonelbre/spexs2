package filters

import (
	. "spexs"
)

type Func func(p *Query, ref *Database) bool

type Conf map[string]interface{}
type CreateFunc func(Conf, Setup) (Func, error)

type Desc struct {
	Name   string
	Desc   string
	Create CreateFunc
}

func Get(name string) (*Desc, bool) {
	for _, e := range All {
		if e.Name == name {
			return &e, true
		}
	}

	f, valid := getFeatureFilter(name)
	return f, valid
}
