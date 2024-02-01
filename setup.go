package main

import (
	"fmt"
	"io"
	"log"
	"regexp"

	"github.com/rogpeppe/rjson"

	"github.com/egonelbre/spexs2/search"
	"github.com/egonelbre/spexs2/search/extenders"
	"github.com/egonelbre/spexs2/search/filters"
	"github.com/egonelbre/spexs2/search/pool"
)

type Printer func(io.Writer, search.Pooler)

type AppSetup struct {
	search.Setup

	conf *Conf

	Order []search.Feature

	Printer    Printer
	printQuery func(io.Writer, *search.Query)

	Features map[string]search.Feature
	Dataset  *Dataset
}

func NewAppSetup(conf *Conf) *AppSetup {
	s := &AppSetup{}
	s.conf = conf

	s.Db, s.Dataset = CreateDatabase(conf)

	s.Order = make([]search.Feature, 0)
	s.Features = make(map[string]search.Feature)

	s.initInput()
	s.initOrder()
	s.initOutput()

	s.initExtender()
	s.initFilters()
	s.initPrinter()

	features := s.Features
	s.PreProcess = func(q *search.Query) error {
		q.CacheValues()
		for _, fn := range features {
			fn(q)
		}
		return nil
	}
	s.PostProcess = func(q *search.Query) error {
		return nil
	}

	return s
}

func (s *AppSetup) initInput() {
	info("init input")
	s.In = pool.NewStack()
}

func (s *AppSetup) initOutput() {
	info("init output")
	size := s.conf.Output.Count
	s.Out = pool.NewPriority(s.Order, size)
}

func (s *AppSetup) initExtender() {
	info("init extender")

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

func (s *AppSetup) makeFilter(name string, data rjson.RawMessage) (search.Filter, error) {
	info("make filter " + name)
	bytes, _ := data.MarshalJSON()

	if isDisabled(bytes) {
		return nil, fmt.Errorf("filter is disabled")
	}

	regRemoveParens := regexp.MustCompile(`\(.*\)`)
	filterName := regRemoveParens.ReplaceAllString(name, "")
	createFilter, ok := filters.Get(filterName)
	if ok {
		return createFilter(&s.Setup, bytes), nil
	}

	// didn't find filter, let's create it from feature
	feature := s.makeFeature(name)
	filter := filters.FromFeature(feature, bytes)
	return filter, nil
}

func (s *AppSetup) makeFilters(conf map[string]rjson.RawMessage) search.Filter {
	info("make filters")
	fns := make([]search.Filter, 0)
	for name, data := range conf {
		fn, err := s.makeFilter(name, data)
		if err == nil {
			fns = append(fns, fn)
		}
	}
	return filters.Compose(fns...)
}

func (s *AppSetup) initFilters() {
	info("init filters")
	s.Extendable = s.makeFilters(s.conf.Extension.Extendable)
	s.Outputtable = s.makeFilters(s.conf.Extension.Outputtable)
}

func (s *AppSetup) initOrder() {
	for _, order := range s.conf.Output.SortBy {
		fn := s.makeFeature(order)
		s.Order = append(s.Order, fn)
	}
}
