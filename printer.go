package main

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"text/template"

	"github.com/egonelbre/spexs2/search"
)

type strFeature func(*search.Query) string

func (s *AppSetup) initPrinter() {
	format := s.conf.Printer.Format
	header := s.conf.Printer.Header
	showHeader := s.conf.Printer.ShowHeader

	if header == "" {
		regHdr, _ := regexp.Compile(`[\{\}]`)
		header = regHdr.ReplaceAllString(format, "")
	}

	features := make(map[string]strFeature)
	featureIdx := 0

	regFeature, _ := regexp.Compile(`[a-zA-Z?() @~,]+`)
	format = regFeature.ReplaceAllStringFunc(format,
		func(call string) string {
			feature, info := s.makeFeatureEx(call)

			var feat strFeature
			if !info {
				feat = func(q *search.Query) string {
					val, _ := feature(q)
					return fmt.Sprintf("%v", val)
				}
			} else {
				feat = func(q *search.Query) string {
					_, info := feature(q)
					return info
				}
			}

			name := fmt.Sprintf("f%v", featureIdx)
			featureIdx++
			features[name] = feat

			return "{{." + name + "}}"
		})

	tmpl, err := template.New("").Parse(format)
	if err != nil {
		log.Println("Unable to create template based on output format. ", format)
		log.Fatal(err)
	}

	printQuery := func(out io.Writer, q *search.Query) {
		values := make(map[string]string)
		for name, fn := range features {
			values[name] = fn(q)
		}

		err = tmpl.Execute(out, values)
		if err != nil {
			log.Println("Unable to output pattern.")
			log.Fatal(err)
		}
	}

	s.printQuery = printQuery

	s.Printer = func(out io.Writer, pool search.Pooler) {
		values := pool.Values()

		if showHeader {
			fmt.Fprint(out, header)
		}

		if !s.conf.Printer.Reverse {
			for _, q := range values {
				printQuery(out, q)
			}
		} else {
			for i := len(values) - 1; i >= 0; i-- {
				q := values[i]
				printQuery(out, q)
			}
		}
	}
}
