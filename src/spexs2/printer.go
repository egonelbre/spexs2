package main

import (
	"fmt"
	"io"
	"log"
	"regexp"
	. "spexs"
	"text/template"
)

type strFeature func(*Query) string

func (s *AppSetup) initPrinter() {
	format := s.conf.Print.Format
	header := s.conf.Print.Header
	showHeader := s.conf.Print.ShowHeader

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
				feat = func(q *Query) string {
					val, _ := q.Memoized(feature)
					return fmt.Sprintf("%v", val)
				}
			} else {
				feat = func(q *Query) string {
					_, info := q.Memoized(feature)
					return info
				}
			}

			name := fmt.Sprintf("f%v", featureIdx)
			featureIdx += 1
			features[name] = feat

			return "{{." + name + "}}"
		})

	tmpl, err := template.New("").Parse(format)
	if err != nil {
		log.Println("Unable to create template based on output format. ", format)
		log.Fatal(err)
	}

	s.Printer = func(out io.Writer, q *Query) {
		if q == nil {
			if showHeader {
				fmt.Print(header)
			}
			return
		}

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
}
