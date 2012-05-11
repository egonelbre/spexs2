package main

import (
	"log"
	. "spexs"
)

type extenderConf map[string]interface{}
type extenderCreator func(conf extenderConf, setup AppSetup) ExtenderFunc

type extenderDesc struct {
	name string
	desc string
	fun  extenderCreator
}

var Extenders = [...]extenderDesc{
	{"simple",
		"extends using the alphabet",
		simpleExtender(SimpleExtender)},
	{"group",
		"extends using the alphabet and group symbols",
		simpleExtender(GroupExtender)},
	{"star",
		"extends using the alphabet and star extension",
		simpleExtender(StarExtender)},
	{"regexp",
		"extends using the alphabet, group symbols and stars",
		simpleExtender(GroupStarExtender)},
}

func getExtender(name string) (extenderCreator, bool) {
	for _, e := range Extenders {
		if e.name == name {
			return e.fun, true
		}
	}
	return nil, false
}

func simpleExtender(f ExtenderFunc) extenderCreator {
	return func(conf extenderConf, setup AppSetup) ExtenderFunc {
		return f
	}
}

func CreateExtender(conf Conf, setup AppSetup) ExtenderFunc {
	if conf.Extension.Method == "" {
		log.Fatal("Extender not defined!")
	}

	create, valid := getExtender(conf.Extension.Method)
	if !valid {
		log.Fatal("No extender named: ", conf.Extension.Method)
	}

	args := conf.Extension.Args[conf.Extension.Method]

	return create(args, setup)
}
