package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"flag"
)

func PrintHelp(conf Conf) {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "spxs [FLAGS] [OPTIONS]\n\n")
	fmt.Fprintf(os.Stderr, "FLAGS: \n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nOPTIONS: \n")

	keys := make([]string, len(conf.Aliases))
	i := 0
	for key := range conf.Aliases {
		keys[i] = key
		i += 1
	}
	sort.Strings(keys)

	for _, name := range keys {
		args := conf.Aliases[name]
		fmt.Fprintf(os.Stderr, "  %s : %s\n", name, args.Desc)
	}

	fmt.Fprintf(os.Stderr, "\nFITNESSES: \n")
	PrintFitnesses(os.Stderr)
	fmt.Fprintf(os.Stderr, "\nEXTENDERS: \n")
	PrintExtenders(os.Stderr)
	fmt.Fprintf(os.Stderr, "\nFILTERS: \n")
	PrintFilters(os.Stderr)
	fmt.Fprintf(os.Stderr, "\nEXAMPLES: \n")
	fmt.Fprintf(os.Stderr, "  spxs -procs=4 inp=data.dna ref=random.dna\n")

	fmt.Fprintf(os.Stderr, "\n")
}

func PrintVersion(conf Conf) {
	fmt.Fprintf(os.Stderr, "%v\n", SPXS_VERSION)
}

func PrintFitnesses(out io.Writer) {
	for name, _ := range Fitnesses {
		fmt.Fprintf(out, "  %s", name)
	}
}

func PrintFilters(out io.Writer) {
	for name, _ := range Filters {
		fmt.Fprintf(out, "  %s", name)
	}
}

func PrintExtenders(out io.Writer) {
	for name, _ := range Extenders {
		fmt.Fprintf(out, "  %s", name)
	}
}