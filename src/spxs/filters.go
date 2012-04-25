package main

import (
	"log"
	"math"
	. "spexs/trie"
)

type filterConf map[string]interface{}
type filterCreator func(filterConf) FilterFunc

type floatFilterConf struct{ Min, Max float64 }

func trueFilter(p *Pattern, ref *Reference) bool {
	return true
}

var Filters = map[string]filterCreator{}

func CreateFilter(conf map[string]map[string]interface{}, setup AppSetup) FilterFunc {
	filters := make([]FilterFunc, 0)

	for name, args := range conf {
		if _, valid := Filters[name]; !valid {
			log.Fatal("No filter named: ", name)
		}
		f := Filters[name](args)
		filters = append(filters, f)
	}

	if len(filters) == 0 {
		return trueFilter
	} else if len(filters) == 1 {
		return filters[0]
	}

	// create a composite filter
	return func(p *Pattern, ref *Reference) bool {
		for _, f := range filters {
			if !f(p, ref) {
				return false
			}
		}
		return true
	}
}

func makeFloatFilter(feature FeatureFunc, config interface{}) FilterFunc {
	var conf floatFilterConf
	conf.Min = math.NaN()
	conf.Max = math.NaN()

	ApplyObject(&config, &conf)

	min, max := conf.Min, conf.Max
	low, high := !math.IsNaN(conf.Min), !math.IsNaN(conf.Max)

	if low && high {
		return func(p *Pattern, ref *Reference) bool {
			return (feature(p, ref) <= max) && (feature(p, ref) >= min)
		}
	} else if low {
		return func(p *Pattern, ref *Reference) bool {
			return feature(p, ref) >= min
		}
	} else if high {
		return func(p *Pattern, ref *Reference) bool {
			return feature(p, ref) <= max
		}
	}

	log.Fatal("Neither min or max was defined for filter.")
	return trueFilter
}

func makeFilter(f FeatureFunc) filterCreator {
	return func(conf filterConf) FilterFunc {
		return makeFloatFilter(f, conf)
	}
}

func initFilters() {
	for name, f := range Features {
		Filters[name] = makeFilter(f.Func)
	}
}
