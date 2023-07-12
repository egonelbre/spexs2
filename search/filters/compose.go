package filters

import (
	"github.com/egonelbre/spexs2/search"
)

func trueFilter(q *search.Query) bool { return true }

func Compose(filters ...search.Filter) search.Filter {
	if len(filters) == 0 {
		return trueFilter
	} else if len(filters) == 1 {
		return filters[0]
	}

	return func(q *search.Query) bool {
		for _, filter := range filters {
			if !filter(q) {
				return false
			}
		}
		return true
	}
}
