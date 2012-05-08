package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	. "spexs"
	"text/template"
)

func CreatePrinter(conf Conf, setup AppSetup) PrinterFunc {
	format := conf.Output.Format

	regExtract, _ := regexp.Compile(`\{([a-zA-Z0-9\-]+)\}`)
	regFixName, _ := regexp.Compile(`-`)

	fixedNames := make(map[string]string)

	formatStrs := regExtract.FindAllStringSubmatch(format, -1)
	for _, tokens := range formatStrs {
		name := tokens[1]
		_, valid := Features[name]
		_, validStr := StrFeatures[name]
		if !(valid || validStr) {
			log.Fatal(errors.New("No valid format parameter: " + name))
		}

		fixedNames[name] = regFixName.ReplaceAllString(name, "")
	}

	format = regExtract.ReplaceAllString(format, `{{.$1}}`)
	format = regFixName.ReplaceAllString(format, "")

	tmpl, err := template.New("").Parse(format)
	if err != nil {
		log.Println("Unable to create template based on output format. ", format)
		log.Fatal(err)
	}

	f := func(out io.Writer, pat *Pattern, ref *Reference) {
		if pat == nil {
			fmt.Print(conf.Output.Format)
			return
		}

		values := make(map[string]interface{})

		for name, fixName := range fixedNames {
			f, valid := Features[name]
			if valid {
				values[fixName] = f.Func(pat, ref)
			}
			fstr, valid := StrFeatures[name]
			if valid {
				values[fixName] = fstr.Func(pat, ref)
			}
		}

		err = tmpl.Execute(out, values)
		if err != nil {
			log.Println("Unable to output pattern.")
			log.Fatal(err)
		}
	}

	//TODO: test printer

	return f
}
