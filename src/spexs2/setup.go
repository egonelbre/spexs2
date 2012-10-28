package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"regexp"
	. "spexs"

	"spexs2/extenders"
	"spexs2/features"
	"spexs2/filters"
	"spexs2/pool"
)

type Printer func(io.Writer, *Query)

type AppSetup struct {
	Setup

	conf *Conf

	Order   Feature
	Printer Printer

	Features map[string]Feature
	Dataset  *Dataset
}

func NewAppSetup(conf *Conf) *AppSetup {
	s := &AppSetup{}
	s.conf = conf

	s.Db, s.Dataset = CreateDatabase(conf)

	s.Features = make(map[string]Feature)

	s.initInput()
	s.initOrder()
	s.initOutput()

	s.initExtender()
	s.initFilters()
	s.initPrinter()

	features := s.Features
	s.PostProcess = func(q *Query) error {
		q.CacheValues()
		for _, fn := range features {
			q.Memoized(fn)
		}
		return nil
	}

	return s
}

func (s *AppSetup) initInput() {
	s.In = pool.NewLifo()
}

func (s *AppSetup) initOutput() {
	if s.conf.Output.Queue == "lifo" {
		s.Out = pool.NewLifo()
		print("lifo")
		return
	}
	size := s.conf.Print.Count
	asc := s.conf.Output.Sort == "asc"
	s.Out = pool.NewPriority(s.Order, size, asc)
}

func (s *AppSetup) initExtender() {
	if s.conf.Extension.Method == "" {
		log.Fatal("Extender not defined!")
	}

	method := s.conf.Extension.Method
	extender, ok := extenders.Get(method)
	if !ok {
		log.Fatal("No extender named: ", method)
	}

	/*args := conf.Extension.Args[conf.Extension.Method]

	ext, err := extender.Create(args)
	if err != nil {
		log.Fatal(err)
	}*/

	s.Extender = extender
}

func parseCall(call string) (name string, groups []string, info bool) {
	regNameArgs, _ := regexp.Compile(`([a-zA-Z0-9]+)(\??)\((.*)\)`)
	regArgs, _ := regexp.Compile("~?[@a-zA-Z0-9]+")
	tokens := regNameArgs.FindStringSubmatch(call)
	if tokens == nil {
		log.Fatalf("Invalid name: %v", call)
	}
	name = tokens[1]
	info = (tokens[2] == "?")
	groups = regArgs.FindAllString(tokens[3], -1)
	return
}

func (s *AppSetup) groupToIds(group string) []int {
	if group == "@" {
		ids := make([]int, len(s.Dataset.Files))
		for i, _ := range ids {
			ids[i] = i
		}
		return ids
	}

	return s.Dataset.Groups[group]
}

func (s *AppSetup) parseFeature(call string) (name string, args [][]int, info bool) {
	name, groups, info := parseCall(call)
	args = make([][]int, len(groups))
	for i, group := range groups {
		args[i] = s.groupToIds(group)
	}
	return
}

func (s *AppSetup) makeFeature(call string) Feature {
	feature, _ := s.makeFeatureEx(call)
	return feature
}

func (s *AppSetup) makeFeatureEx(call string) (Feature, bool) {
	name, args, info := s.parseFeature(call)

	normalized := fmt.Sprintf("%+v%+v", name, args)
	if feature, ok := s.Features[normalized]; ok {
		return feature, info
	}

	create, ok := features.Get(name)
	if !ok {
		log.Fatal("No feature named ", name)
	}

	feature, err := features.CallCreateWithArgs(create, args)
	if err != nil {
		log.Fatal(err)
	}

	s.Features[normalized] = feature
	return feature, info
}

func (s *AppSetup) makeFilter(name string, data json.RawMessage) Filter {
	bytes, _ := data.MarshalJSON()

	createFilter, ok := filters.Get(name)
	if ok {
		return createFilter(s.Setup, bytes)
	}

	// didn't find filter, let's create it from feature
	feature := s.makeFeature(name)
	filter := filters.FeatureFilter(feature, bytes)
	return filter
}

func (s *AppSetup) makeFilters(conf map[string]json.RawMessage) Filter {
	fns := make([]Filter, 0)
	for name, data := range conf {
		fn := s.makeFilter(name, data)
		fns = append(fns, fn)
	}
	return filters.Compose(fns)
}

func (s *AppSetup) initFilters() {
	s.Extendable = s.makeFilters(s.conf.Extension.Filter)
	s.Outputtable = s.makeFilters(s.conf.Output.Filter)
}

func lengthFitness(q *Query) float64 {
	return 1 / float64(q.Len())
}

func (s *AppSetup) initOrder() {
	order := s.conf.Output.Order
	if order == "" {
		log.Fatal("Output ordering not defined!")
	}
	s.Order = s.makeFeature(order)
}
