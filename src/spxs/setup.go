package main

import (
	. "spexs/trie"
)

const MAX_POOL_SIZE = 1024 * 1024 * 1024

type TrieFitnessCreator func(interface{}) FitnessFunc
type TrieExtenderCreator func(interface{}) ExtenderFunc

func lengthFitness(p *Pattern) float64 {
	return 1 / float64(p.Len())
}

type PatternFilterCreator func(limit int) FilterFunc

func CreateInput(conf Conf, setup Setup) Pooler {
	in := NewPriorityPool(lengthFitness, MAX_POOL_SIZE)
	in.Put(NewFullPattern(setup.Ref))
	return in
}

func CreateOutput(conf Conf, setup Setup) Pooler {
	size := conf.Output.Count
	if size < 0 {
		size = MAX_POOL_SIZE
	}
	return NewPriorityPool(setup.Fitness, size)
}

func CreateSetup(conf Conf) Setup {
	var setup Setup
	setup.Ref = CreateReference(conf)

	setup.In = CreateInput(conf, setup)
	setup.Fitness = CreateFitness(conf, setup)
	setup.Out = CreateOutput(conf, setup)

	setup.Extender = CreateExtender(conf, setup)
	setup.Extendable = CreateFilter(conf.Extension.Filter, setup)
	setup.Outputtable = CreateFilter(conf.Output.Filter, setup)

	return setup
}
