package filters

import (
	"encoding/json"
	"log"
	"math"

	. "github.com/egonelbre/spexs2/search"
)

type minmax struct{ Min, Max float64 }

type featureFilter struct {
	Feature Feature
	Min     float64
	Max     float64
}

func NewFeatureFilter(f Feature, data []byte) Filter {
	var conf struct{ Min, Max float64 }
	conf.Min = math.NaN()
	conf.Max = math.NaN()

	err := json.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal(err)
	}

	min, max := conf.Min, conf.Max
	if math.IsNaN(conf.Min) {
		min = math.Inf(-1)
	}
	if math.IsNaN(conf.Max) {
		max = math.Inf(1)
	}

	return &featureFilter{f, min, max}
}

func (f *featureFilter) Accepts(q *Query) bool {
	val, _ := f.Feature.Evaluate(q)
	return (f.Min <= val) && (val <= f.Max)
}
