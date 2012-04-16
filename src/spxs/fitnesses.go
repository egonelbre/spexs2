package main

import (
	"log"
	. "spexs/trie"
)

type fitnessConf map[string]interface{}

type fitnessCreator func(fitnessConf, Setup) FitnessFunc

func simpleFitness(f FitnessFunc) fitnessCreator {
	return func(conf fitnessConf, setup Setup) FitnessFunc {
		return f
	}
}

var fitnesses = map[string]fitnessCreator{
	"def": simpleFitness(
		func(p *Pattern) float64 {
			return float64(p.Len() * p.Pos.Len())
		}),
	"len": simpleFitness(
		func(p *Pattern) float64 {
			return float64(p.Len())
		}),
	"count": simpleFitness(
		func(p *Pattern) float64 {
			return float64(p.Pos.Len())
		}),
	"complexity": simpleFitness(
		func(p *Pattern) float64 {
			return float64(p.Complexity())
		}),
	"p-value": func(conf fitnessConf, setup Setup) FitnessFunc {
		return func(p *Pattern) float64 {
			return p.PValue(setup.Ref)
		}
	},
}

func CreateFitness(conf Conf, setup Setup) FitnessFunc {
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
