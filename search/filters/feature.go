package filters

import (
	"encoding/json"
	"log"
	"math"
	
	. "github.com/egonelbre/spexs2/search"
)

type minmax struct{ Min, Max float64 }

func FeatureFilter(feature Feature, data []byte) Filter {
	var conf struct{ Min, Max float64 }
	conf.Min = math.NaN()
	conf.Max = math.NaN()

	err := json.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal(err)
	}

	min, max := conf.Min, conf.Max
	low, high := !math.IsNaN(conf.Min), !math.IsNaN(conf.Max)

	if low && high {
		return func(q *Query) bool {
			val, _ := feature(q)
			return (min <= val) && (val <= max)
		}
	} else if low {
		return func(q *Query) bool {
			val, _ := feature(q)
			return min <= val
		}
	} else if high {
		return func(q *Query) bool {
			val, _ := feature(q)
			return val <= max
		}
	}

	return func(q *Query) bool { return true }
}
