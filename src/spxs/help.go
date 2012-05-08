package main

import (
	"flag"
	"log"
	"os"
	"sort"
)

var lgh = log.New(os.Stderr, "", 0)

func PrintHelp(conf Conf) {
	lgh.Printf("Usage of %s:\n", os.Args[0])
	lgh.Printf("spxs [FLAGS] [OPTIONS]\n\n")
	PrintVersion()
	lgh.Printf("FLAGS: \n")
	flag.PrintDefaults()
	lgh.Printf("\nALIASES: \n")

	keys := make([]string, len(conf.Aliases))
	i := 0
	for key := range conf.Aliases {
		keys[i] = key
		i += 1
	}
	sort.Strings(keys)

	for _, name := range keys {
		args := conf.Aliases[name]
		lgh.Printf("  %s : %s\n", name, args.Desc)
	}

	PrintStrFeatures()
	PrintFeatures()
	PrintFitnesses()
	PrintFilters()
	PrintExtenders()

	lgh.Printf("\nEXAMPLES: \n")
	lgh.Printf("  spxs -procs=4 inp=data.dna ref=random.dna\n")

	lgh.Printf("\n")
}

func PrintVersion() {
	lgh.Printf("%v\n", theVersion)
	lgh.Printf("%v\n", theBuildTime)
}

func PrintStrFeatures() {
	lgh.Printf("\nSTRING FEATURES: \n")
	for name, f := range StrFeatures {
		lgh.Printf("  %s : %s\n", name, f.Desc)
	}
}

func PrintFeatures() {
	lgh.Printf("\nFEATURES: \n")
	for name, f := range Features {
		lgh.Printf("  %s : %s\n", name, f.Desc)
	}
}

func PrintFitnesses() {
	lgh.Printf("\nFITNESSES: \n")
	lgh.Printf("  +[FEATURES]\n")
	for name, _ := range Fitnesses {
		lgh.Printf("  %s\n", name)
	}
}

func PrintFilters() {
	lgh.Printf("\nFILTERS: \n")
	lgh.Printf("  +[FEATURES]\n")
	for name, _ := range Filters {
		lgh.Printf("  %s\n", name)
	}
}

func PrintExtenders() {
	lgh.Printf("\nEXTENDERS: \n")
	for name, _ := range Extenders {
		lgh.Printf("  %s\n", name)
	}
}
