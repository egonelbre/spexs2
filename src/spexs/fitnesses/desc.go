package fitnesses

import (
	. "spexs"
)

type Func func(p *Pattern, ref *Reference) float64

type Conf map[string]interface{}
type CreateFunc func(Conf) (Func, error)

type Desc struct {
	Name string
	Desc string
	Create CreateFunc
}

func Get(name string) (*Desc, bool) {
	for _, e := range All {
		if e.Name == name {
			return &e, true
		}
	}

	return getFeatureFitness(name)
}