package main

import (
	"log"
	. "spexs"
)

var fitnesses = map[string]TrieFitnessFunc{
	"def": func(p *TrieNode) float64 {
		return float64(p.Len() * p.Pos.Len())
	},
	"len": func(p *TrieNode) float64 {
		return float64(p.Len())
	},
	"count": func(p *TrieNode) float64 {
		return float64(p.Pos.Len())
	},
	"complexity": func(p *TrieNode) float64 {
		return float64(p.Complexity())
	},
}

func CreateFitness(conf Conf, setup Setup) TrieFitnessFunc {
	if conf.Output.Order == "" {
		log.Fatal("Output ordering not defined!")
	}

	//TODO: read in additional args

	if _, valid := fitnesses[conf.Output.Order]; !valid {
		log.Fatal("No ordering/fitness named: ", conf.Output.Order)
	}

	return fitnesses[conf.Output.Order]
}
