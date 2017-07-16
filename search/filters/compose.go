package filters

import (
	. "github.com/egonelbre/spexs2/search"
)

func trueFilter(q *Query) bool { return true }

func Compose(filters ...Filter) Filter {
	if len(filters) == 0 {
		return trueFilter
	} else if len(filters) == 1 {
		return filters[0]
	}

	return func(q *Query) bool {
		for _, filter := range filters {
			if !filter(q) {
				return false
			}
		}
		return true
	}
}
