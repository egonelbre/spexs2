package filters

import (
	. "github.com/egonelbre/spexs2/search"
)

type trueFilter struct{}

func (f *trueFilter) Accepts(q *Query) bool {
	return true
}

type compositeFilter []Filter

func (fs compositeFilter) Accepts(q *Query) bool {
	for _, filter := range fs {
		if !filter.Accepts(q) {
			return false
		}
	}
	return true
}

func Compose(filters []Filter) Filter {
	if len(filters) == 0 {
		return &trueFilter{}
	} else if len(filters) == 1 {
		return filters[0]
	}

	return (compositeFilter)(filters)
}
