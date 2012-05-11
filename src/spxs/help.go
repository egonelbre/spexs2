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
	info := make(map[string]string)
	for name, data := range StrFeatures {
		info[name] = data.Desc
	}

	printSection("String Features", info)
}

func PrintFeatures() {
	info := make(map[string]string)
	for name, data := range Features {
		info[name] = data.Desc
	}

	printSection("Features", info)
}

func PrintFitnesses() {
	info := make(map[string]string)
	for name, _ := range Fitnesses {
		info[name] = ""
	}
	info["+[FEATURES]"] = ""

	printSection("Fitnesses", info)
}

func PrintFilters() {
	info := make(map[string]string)
	for name, _ := range Filters {
		info[name] = ""
	}
	info["+[FEATURES]"] = ""

	printSection("Filters", info)
}

func PrintExtenders() {
	printCaption("Extenders")
	for _, e := range Extenders {
		printItem(e.name, e.desc)
	}
}
