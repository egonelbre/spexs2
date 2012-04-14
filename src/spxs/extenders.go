package main

import (
	"log"
	. "spexs"
)

type extenderConf interface{}
type extenderCreator func(conf extenderConf, setup Setup) TrieExtenderFunc

func simpleExtender(f TrieExtenderFunc) extenderCreator {
	return func(conf extenderConf, setup Setup) TrieExtenderFunc {
		return f
	}
}

var extenders = map[string]extenderCreator{
	"simple": simpleExtender(SimpleExtender),
	"group":  simpleExtender(GroupExtender),
	"star":   simpleExtender(StarExtender),
	"regexp": simpleExtender(GroupStarExtender),
}

func CreateExtender(conf Conf, setup Setup) TrieExtenderFunc {
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
