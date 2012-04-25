package main

import (
	"log"
	. "spexs/trie"
)

type fitnessConf map[string]interface{}
type fitnessCreator func(fitnessConf, AppSetup) FitnessFunc

func simpleFitness(f FitnessFunc) fitnessCreator {
	return func(conf fitnessConf, setup AppSetup) FitnessFunc {
		return f
	}
}

var Fitnesses = map[string]fitnessCreator{
}

func CreateFitness(conf Conf, setup AppSetup) FitnessFunc {
	if conf.Output.Order == "" {
		log.Fatal("Output ordering not defined!")
	}

	fitnessCreate, valid := Fitnesses[conf.Output.Order]
	if !valid {
		log.Fatal("No ordering/fitness named: ", conf.Output.Order)
	}
	args := conf.Output.Args[conf.Output.Order]

	return fitnessCreate(args, setup)
}

func makeFloatFitness(f FeatureFunc) fitnessCreator {
	return func(conf fitnessConf, setup AppSetup) FitnessFunc {
		return func(p *Pattern) float64 {
			return f(p, setup.Ref)
		}
	}
}

func initFitnesses() {
	for name, f := range Features {
		Fitnesses[name] = makeFloatFitness(f)
	}
}
