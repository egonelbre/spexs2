package main

import (
	"log"
	. "spexs"
)

type fitnessConf map[string]interface{}

type fitnessCreator func(fitnessConf, Setup) TrieFitnessFunc

func simpleFitness(f TrieFitnessFunc) fitnessCreator {
	return func(conf fitnessConf, setup Setup) TrieFitnessFunc {
		return f
	}
}

var fitnesses = map[string]fitnessCreator{
	"def": simpleFitness(
		func(p *TrieNode) float64 {
			return float64(p.Len() * p.Pos.Len())
		}),
	"len": simpleFitness(
		func(p *TrieNode) float64 {
			return float64(p.Len())
		}),
	"count": simpleFitness(
		func(p *TrieNode) float64 {
			return float64(p.Pos.Len())
		}),
	"complexity": simpleFitness(
		func(p *TrieNode) float64 {
			return float64(p.Complexity())
		}),
	"p-value": func(conf fitnessConf, setup Setup) TrieFitnessFunc {
		return func(p *TrieNode) float64 {
			return p.PValue(setup.Ref)
		}
	},
}

func CreateFitness(conf Conf, setup Setup) TrieFitnessFunc {
	if conf.Output.Order == "" {
		log.Fatal("Output ordering not defined!")
	}

	fitnessCreate, valid := fitnesses[conf.Output.Order]
	if !valid {
		log.Fatal("No ordering/fitness named: ", conf.Output.Order)
	}

	args := conf.Output.Args[conf.Output.Order]

	return fitnessCreate(args, setup)
}
