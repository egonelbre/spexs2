package main

import (
	"log"
	"math"
	. "spexs"
)

type filterConf map[string]interface{}
type filterCreator func(filterConf, Setup) TrieFilterFunc

type genericFilterConf struct{ Min, Max float64 }

func trueFilter(p *TrieNode) bool {
	return true
}

var Filters = map[string] filterCreator {
	"length":  func(conf filterConf, setup Setup) TrieFilterFunc {
		return genericFilter(func(p *TrieNode) float64 {
			return float64(p.Len())
	}, conf)},
	"count": func(conf filterConf, setup Setup) TrieFilterFunc {
		return genericFilter(func(p *TrieNode) float64 {
			return float64(p.Pos.Len()) / float64(len(setup.Ref.Pats))
	}, conf)},
	"p-value": func(conf filterConf, setup Setup) TrieFilterFunc {
		return genericFilter(func(p *TrieNode) float64 {
			return p.PValue(setup.Ref)
	}, conf)},
}

func CreateFilter(conf map[string]map[string]interface{}, setup Setup) TrieFilterFunc {
	filters := make([]TrieFilterFunc, 0)

	for name, args := range conf {
		if _, valid := Filters[name]; !valid {
			log.Fatal("No filter named: ", name)
		}
		f := Filters[name](args, setup)
		filters = append(filters, f)
	}

	if len(filters) == 0 {
		return trueFilter
	} else if len(filters) == 1 {
		return filters[0]
	}

	// create a composite filter
	return func(p *TrieNode) bool {
		for _, f := range filters {
			if !f(p) {
				return false
			}
		}
		return true
	}
}

type valueFunc func(p *TrieNode) float64

func genericFilter(value valueFunc, config interface{}) TrieFilterFunc {
	var conf genericFilterConf
	conf.Min = math.NaN()
	conf.Max = math.NaN()

	ApplyObject(&config, &conf)

	min, max := conf.Min, conf.Max
	low, high := !math.IsNaN(conf.Min), !math.IsNaN(conf.Max)

	if low && high {
		return func(p *TrieNode) bool {
			return (value(p) <= max) && (value(p) > min)
		}
	} else if low {
		return func(p *TrieNode) bool {
			return value(p) > min
		}
	} else if high {
		return func(p *TrieNode) bool {
			return value(p) <= max
		}
	}

	log.Fatal("Neither min or max was defined for count filter.")	
	return trueFilter
}
