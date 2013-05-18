package features

import (
	. "spexs"
	"stats/binom"
)

// binomial p-value
func Binom(fore, back []int) Feature {
	return func(q *Query) (float64, string) {
		totalFore := count(q.Db.Total, fore) + 1
		totalBack := count(q.Db.Total, back) + 1

		matches := q.Matches()
		countFore := count(matches, fore)
		countBack := count(matches, back)

		p := float64(countBack) / float64(totalBack)

		return binom.P(countFore, totalFore, p), ""
	}
}

