package main

import (
	"log"
	. "spexs"
)

var extenders = map[string] TrieExtenderFunc {
	"simple": SimpleExtender,
	"group":  GroupExtender,
	"star":   StarExtender,
	"regexp": GroupStarExtender,
}

func CreateExtender(conf Conf, setup Setup) TrieExtenderFunc {
	if conf.Extension.Method == "" {
		log.Fatal("Extender not defined!")
	}

	//TODO: read in additional args

	if _, valid := extenders[conf.Extension.Method]; !valid {
		log.Fatal("No extender named: ", conf.Extension.Method)
	}
	
	return extenders[conf.Extension.Method]
}