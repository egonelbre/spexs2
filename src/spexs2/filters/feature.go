package feature

import (
	"encoding/json"
	"math"
	. "spexs"
)

type minmax struct{ Min, Max float64 }

func FeatureFilter(feature FeatureFunc, data []byte) FilterFunc {
	var conf struct{ Min, Max float64 }
	conf.Min = math.NaN()
	conf.Max = math.NaN()

	err := json.Unmarshal(data, &conf)
	if err != nil {
		panic(err)
	}

	min, max := conf.Min, conf.Max
	low, high := !math.IsNaN(conf.Min), !math.IsNaN(conf.Max)

	if low && high {
		return func(q *Query) bool {
			val := p.Memoized(feature)
			return (min <= val) && (val <= max)
		}
	} else if low {
		return func(q *Query) bool {
			val := p.Memoized(feature)
			return (min <= val) && (val <= max)
		}
	} else if high {
		return func(q *Query) bool {
			val := p.Memoized(feature)
			return (min <= val) && (val <= max)
		}
	}

	return func(q *Query) bool { return true }
}
