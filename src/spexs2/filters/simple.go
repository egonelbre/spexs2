package filters

import (
	. "spexs"
)

// don't allow start or end to be grouping token
func NoGroupingEnds(conf Conf, setup Setup) FilterFunc {
	return func(q *Query) bool {
		start := q.Pat[0]
		end := q.Pat[len(q.Pat)-1]
		return !start.IsGroup && !start.IsStar && !end.IsGroup
	}
}
