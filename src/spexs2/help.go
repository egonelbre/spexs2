package main

import (
	"flag"
	"log"
	"os"

	"spexs/extenders"
	"spexs/features"
	"spexs/filters"
)

var lgh = log.New(os.Stderr, "", 0)

func PrintHelp() {
	lgh.Printf("Usage of %s:\n", os.Args[0])
	lgh.Printf("%s [FLAGS] [OPTIONS]\n\n", os.Args[0])
	lgh.Printf("Flags: \n")
	flag.PrintDefaults()

	lgh.Printf("\n### Extenders\n")
	lgh.Printf(extenders.Help)

	lgh.Printf("\n### Filters\n")
	lgh.Printf(filters.Help)

	lgh.Printf("\n### Features\n")
	lgh.Printf(features.Help)

	//lgh.Printf("\n### Examples\n")
	//lgh.Printf("  spexs2 -conf=conf.json inp=data.dna ref=random.dna\n")

	lgh.Printf("\n")
}

func PrintVersion() {
	lgh.Printf("%v\n", theVersion)
	lgh.Printf("%v\n", theBuildTime)
}