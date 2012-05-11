package extenders

import (
	. "spexs"
)

type Func func(p *Pattern, ref *Reference) Patterns

type Conf map[string]interface{}
type CreateFunc func(conf Conf) (Func, error)

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
	return nil, false
}

func wrap(f Func) CreateFunc {
	return func(conf Conf) (Func, error) {
		return f, nil
	}
}
