package main

import (
	"flag"
	"log"
	"os"

	"github.com/egonelbre/spexs2/search/extenders"
	"github.com/egonelbre/spexs2/search/features"
	"github.com/egonelbre/spexs2/search/filters"
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
	lgh.Printf("%v\n", buildversion)
	lgh.Printf("%v\n", buildtime)
}
