package main

import (
	"fmt"
	"io"
	"launchpad.net/rjson"
	"log"
	"regexp"

	. "github.com/egonelbre/spexs2/search"
	"github.com/egonelbre/spexs2/search/extenders"
	"github.com/egonelbre/spexs2/search/filters"
	"github.com/egonelbre/spexs2/search/pool"
)

type Printer func(io.Writer, Pooler)

type AppSetup struct {
	Setup

	conf *Conf

	Order []Feature

	Printer    Printer
	printQuery func(io.Writer, *Query)

	Features map[string]Feature
	Dataset  *Dataset
}

func NewAppSetup(conf *Conf) *AppSetup {
	s := &AppSetup{}
	s.conf = conf

	s.Db, s.Dataset = CreateDatabase(conf)

	s.Order = make([]Feature, 0)
	s.Features = make(map[string]Feature)

	s.initInput()
	s.initOrder()
	s.initOutput()

	s.initExtender()
	s.initFilters()
	s.initPrinter()

	features := s.Features
	s.PreProcess = func(q *Query) error {
		q.CacheValues()
		for _, fn := range features {
			fn(q)
		}
		return nil
	}
	s.PostProcess = func(q *Query) error {
		q.Loc = nil
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

func (s *AppSetup) makeFilter(name string, data rjson.RawMessage) (Filter, error) {
	info("make filter " + name)
	bytes, _ := data.MarshalJSON()

	if isDisabled(bytes) {
		return nil, fmt.Errorf("filter is disabled")
	}

	regRemoveParens, _ := regexp.Compile(`\(.*\)`)
	filterName := regRemoveParens.ReplaceAllString(name, "")
	createFilter, ok := filters.Get(filterName)
	if ok {
		return createFilter(s.Setup, bytes), nil
	}

	// didn't find filter, let's create it from feature
	feature := s.makeFeature(name)
	filter := filters.FeatureFilter(feature, bytes)
	return filter, nil
}

func (s *AppSetup) makeFilters(conf map[string]rjson.RawMessage) Filter {
	info("make filters")
	fns := make([]Filter, 0)
	for name, data := range conf {
		fn, err := s.makeFilter(name, data)
		if err == nil {
			fns = append(fns, fn)
		}
	}
	return filters.Compose(fns)
}

func (s *AppSetup) initFilters() {
	info("init filters")
	s.Extendable = s.makeFilters(s.conf.Extension.Extendable)
	s.Outputtable = s.makeFilters(s.conf.Extension.Outputtable)
}

func lengthFitness(q *Query) float64 {
	return 1 / float64(q.Len())
}

func (s *AppSetup) initOrder() {
	for _, order := range s.conf.Output.SortBy {
		fn := s.makeFeature(order)
		s.Order = append(s.Order, fn)
	}
}
