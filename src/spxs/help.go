package main

import (
	"flag"
	"log"
	"os"
	"sort"

	"spexs/filters"
	"spexs/extenders"
	"spexs/features"
	"spexs/fitnesses"
)

var lgh = log.New(os.Stderr, "", 0)

func PrintHelp(conf Conf) {
	lgh.Printf("Usage of %s:\n", os.Args[0])
	lgh.Printf("spxs [FLAGS] [OPTIONS]\n\n")
	lgh.Printf("FLAGS: \n")
	flag.PrintDefaults()

	PrintAliases(conf)
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

func printCaption(caption string) {
	lgh.Printf("\n%s: \n", caption)
}

func printItem(name, info string) {
	if info != "" {
		lgh.Printf("  %s : %s\n", name, info)
	} else {
		lgh.Printf("  %s\n", name)
	}
}

func printSection(caption string, data map[string]string) {

	i := 0
	names := make([]string, len(data))
	for name := range data {
		names[i] = name
		i += 1
	}
	sort.Strings(names)

	printCaption(caption)
	for _, name := range names {
		printItem(name, data[name])
	}
}

func PrintAliases(conf Conf) {
	info := make(map[string]string)
	for name, data := range conf.Aliases {
		info[name] = data.Desc
	}

	printSection("Aliases", info)
}

func PrintStrFeatures() {
	printCaption("Features")
	for _, e := range features.Str {
		printItem(e.Name, e.Desc)
	}
}

func PrintFeatures() {
	printCaption("Features")
	for _, e := range features.All {
		printItem(e.Name, e.Desc)
	}
}

func PrintFitnesses() {
	printCaption("Fitnesses")
	printItem("+[FEATURES]", "")

	for _, e := range fitnesses.All {
		printItem(e.Name, e.Desc)
	}
}

func PrintFilters() {
	printCaption("Filters")
	printItem("+[FEATURES]", "")

	for _, e := range filters.All {
		printItem(e.Name, e.Desc)
	}
}

func PrintExtenders() {
	printCaption("Extenders")
	for _, e := range extenders.All {
		printItem(e.Name, e.Desc)
	}
}
