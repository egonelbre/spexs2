package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	. "spexs"
	"text/template"

	"spexs/features"
)

func CreatePrinter(conf Conf, setup AppSetup) PrinterFunc {

	format := conf.Output.Format
	header := conf.Output.Header
	if header == "" {
		regHdr, _ := regexp.Compile(`[\{\}]`)
		header = regHdr.ReplaceAllString(format, "")
	}

	regExtract, _ := regexp.Compile(`\{([a-zA-Z0-9\-]+)\}`)
	regFixName, _ := regexp.Compile(`-`)

	fixedNames := make(map[string]string)

	feats := make(map[string]features.Func)
	strFeats := make(map[string]features.StrFunc)

	formatStrs := regExtract.FindAllStringSubmatch(format, -1)
	for _, tokens := range formatStrs {
		name := tokens[1]

		f, valid := features.Get(name)
		fs, validStr := features.GetStr(name)

		if !(valid || validStr) {
			log.Fatal(errors.New("No valid format parameter: " + name))
		}

		if valid {
			feats[name] = f.Func
		}

		if validStr {
			strFeats[name] = fs.Func
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
			fmt.Print(header)
			return
		}

		values := make(map[string]interface{})

		for name, fixName := range fixedNames {
			f, valid := feats[name]
			if valid {
				values[fixName] = f(pat, ref)
			}
			fstr, valid := strFeats[name]
			if valid {
				values[fixName] = fstr(pat, ref)
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
