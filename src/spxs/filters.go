package main

import (
	"log"
	"math"
	. "spexs/trie"
)

type filterConf map[string]interface{}
type filterCreator func(filterConf) FilterFunc

type genericFilterConf struct{ Min, Max float64 }

func trueFilter(p *Pattern, ref *Reference) bool {
	return true
}

var Filters = map[string]filterCreator{
	"length": func(conf filterConf) FilterFunc {
		return genericFilter(func(p *Pattern, ref *Reference) float64 {
			return float64(p.Len())
		}, conf)
	},
	"count": func(conf filterConf) FilterFunc {
		return genericFilter(func(p *Pattern, ref *Reference) float64 {
			return float64(p.Pos.Len()) / float64(len(ref.Pats))
		}, conf)
	},
	"p-value": func(conf filterConf) FilterFunc {
		return genericFilter(func(p *Pattern, ref *Reference) float64 {
			return p.PValue(ref)
		}, conf)
	},
}

func CreateFilter(conf map[string]map[string]interface{}, setup Setup) FilterFunc {
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

type valueFunc func(*Pattern, *Reference) float64

func genericFilter(value valueFunc, config interface{}) FilterFunc {
	var conf genericFilterConf
	conf.Min = math.NaN()
	conf.Max = math.NaN()

	ApplyObject(&config, &conf)

	min, max := conf.Min, conf.Max
	low, high := !math.IsNaN(conf.Min), !math.IsNaN(conf.Max)

	if low && high {
		return func(p *Pattern, ref *Reference) bool {
			return (value(p, ref) <= max) && (value(p, ref) > min)
		}
	} else if low {
		return func(p *Pattern, ref *Reference) bool {
			return value(p, ref) > min
		}
	} else if high {
		return func(p *Pattern, ref *Reference) bool {
			return value(p, ref) <= max
		}
	}

	log.Fatal("Neither min or max was defined for filter.")
	return trueFilter
}
