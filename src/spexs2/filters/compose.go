package filters

import (
	. "spexs"
)

func trueFilter(q *Query) bool { return true }

func Compose(filters []FilterFunc) FilterFunc {
	if len(filters) == 0 {
		return trueFilter
	} else if len(filters) == 1 {
		return filters[0]
	}

	return func(q *Query) bool {
		for _, filter := range filters {
			if !filter(p) {
				return false
			}
		}
		return true
	}
}
