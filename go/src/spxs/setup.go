package main

import (
	. "spexs"
)

const MAX_POOL_SIZE = 1024*1024*1024

type Setup struct {
	Ref        *UnicodeReference
	Out        TriePooler
	In         TriePooler

	Extender   TrieExtenderFunc
	
	Extendable TrieFilterFunc
	Outputtable TrieFilterFunc

	Fitness    TrieFitnessFunc
}


type TrieFitnessCreator func(interface{}) TrieFitnessFunc
type TrieExtenderCreator func(interface{}) TrieExtenderFunc

func lengthFitness(p *TrieNode) float64 {
	return 1 / float64(p.Len())
}

type PatternFilterCreator func(limit int) TrieFilterFunc

func CreateInput(conf Conf) TriePooler {
	return NewPriorityPool(lengthFitness, MAX_POOL_SIZE)
}

func CreateOutput(conf Conf, fitness TrieFitnessFunc) TriePooler {
	size := conf.Output.Count
	if size < 0 {
		size = MAX_POOL_SIZE
	}
	return NewPriorityPool(fitness, size)
}

func CreateSetup(conf Conf) Setup {
	var setup Setup

	setup.Ref = CreateReference(conf)
	
	setup.In = CreateInput(conf)
	setup.Fitness = CreateFitness(conf, setup)
	setup.Out = CreateOutput(conf, setup.Fitness)

	setup.Extender = CreateExtender(conf, setup)
	setup.Extendable = CreateFilter(conf.Extension.Filter, setup)
	setup.Outputtable = CreateFilter(conf.Output.Filter, setup)

	return setup
}