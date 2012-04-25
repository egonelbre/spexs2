package main

import (
	"fmt"
	"io"
	"log"
	. "spexs/trie"
	"text/template"
)

type printerArgs struct {
	Str     string
	Regexp  string
	Fitness float64
	Length  int
	Count   int
	PValue  float64
}

func CreatePrinter(conf Conf, setup AppSetup) PrinterFunc {
	tmpl, err := template.New("").Parse(conf.Output.Format)
	if err != nil {
		log.Println("Unable to create template based on output format.")
		log.Fatal(err)
	}

	f := func(out io.Writer, pat *Pattern, ref *Reference) {
		if pat == nil {
			fmt.Println(conf.Output.Format)
			return
		}

		values := make(map[string]interface{})
		
		for name, f := range Features {
			values[name] = f.Func(pat, ref)
		}

		values["str"] = pat.String()
		values["regexp"] = setup.Ref.ReplaceGroups(pat.String())

		err = tmpl.Execute(out, values)
		if err != nil {
			log.Println("Unable to output pattern.")
			log.Fatal(err)
		}
	}

	//TODO: test printer

	return f
}
