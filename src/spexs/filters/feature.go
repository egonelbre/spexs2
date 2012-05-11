package filters

import (
	. "spexs"
	"spexs/features"
	"math"
	"utils"
	"errors"
)

type floatConf struct{ Min, Max float64 }

func makeFilter(feature features.Func, config interface{}) (f Func, err error) {
	err = nil
	var conf floatConf
	conf.Min = math.NaN()
	conf.Max = math.NaN()

	utils.ApplyObject(&config, &conf)

	min, max := conf.Min, conf.Max
	low, high := !math.IsNaN(conf.Min), !math.IsNaN(conf.Max)

	f = nil
	if low && high {
		f = func(p *Pattern, ref *Reference) bool {
			return (feature(p, ref) <= max) && (feature(p, ref) >= min)
		}
	} else if low {
		f = func(p *Pattern, ref *Reference) bool {
			return feature(p, ref) >= min
		}
	} else if high {
		f = func(p *Pattern, ref *Reference) bool {
			return feature(p, ref) <= max
		}
	}

	if f == nil {
		return trueFilter, errors.New("Neither min or max was defined for filter.")
	}

	return f, nil
}

func wrap(f features.Func) CreateFunc {
	return func(conf Conf) (Func, error) {
		return makeFilter(f, conf)
	}
}

func getFeatureFilter(name string) (*Desc, bool) {
	f, valid := features.Get(name)
	if !valid {
		return nil, false
	}
	return &Desc{f.Name, f.Desc, wrap(f.Func)}, true
}
