package fitnesses

import (
	"spexs/features"
)

func wrap(f features.Func) CreateFunc {
	return func(conf Conf) (Func, error) {
		return Func(f), nil
	}
}

func getFeatureFitness(name string) (*Desc, bool) {
	f, valid := features.Get(name)
	if !valid {
		return nil, false
	}
	return &Desc{f.Name, f.Desc, wrap(f.Func)}, true
}
