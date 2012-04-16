package main

import (
	"log"
	. "spexs/trie"
)

type extenderConf interface{}
type extenderCreator func(conf extenderConf, setup AppSetup) ExtenderFunc

func simpleExtender(f ExtenderFunc) extenderCreator {
	return func(conf extenderConf, setup AppSetup) ExtenderFunc {
		return f
	}
}

var extenders = map[string]extenderCreator{
	"simple": simpleExtender(SimpleExtender),
	"group":  simpleExtender(GroupExtender),
	"star":   simpleExtender(StarExtender),
	"regexp": simpleExtender(GroupStarExtender),
}

func CreateExtender(conf Conf, setup AppSetup) ExtenderFunc {
	if conf.Extension.Method == "" {
		log.Fatal("Extender not defined!")
	}

	extenderCreate, valid := extenders[conf.Extension.Method]
	if !valid {
		log.Fatal("No extender named: ", conf.Extension.Method)
	}

	args := conf.Extension.Args[conf.Extension.Method]

	return extenderCreate(args, setup)
}
