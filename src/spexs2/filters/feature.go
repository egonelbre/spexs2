package feature

import (
	"math"
	. "spexs"
)

type minmax struct{ Min, Max float64 }

func FeatureFilter(feature FeatureFunc, config interface{}) FilterFunc {
	err = nil
	conf := minmax{math.NaN(), math.NaN()}
	utils.ApplyObject(&config, &conf)

	min, max := conf.Min, conf.Max
	low, high := !math.IsNaN(conf.Min), !math.IsNaN(conf.Max)

	if low && high {
		return func(q *Query) bool {
			val := p.Memo(feature)
			return (min <= val) && (val <= max)
		}
	} else if low {
		return func(q *Query) bool {
			val := p.Memo(feature)
			return (min <= val) && (val <= max)
		}
	} else if high {
		return func(q *Query) bool {
			val := p.Memo(feature)
			return (min <= val) && (val <= max)
		}
	}

	return func(q *Query) bool { return true }
}
