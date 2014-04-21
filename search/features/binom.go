package features

import (
	. "github.com/egonelbre/spexs2/search"
	"github.com/egonelbre/spexs2/stats/binom"
)

// binomial p-value
func Binom(fore, back []int) Feature {
	return func(q *Query) (float64, string) {
		totalFore := count(q.Db.Total, fore) + 1
		totalBack := count(q.Db.Total, back) + 1

		matches := q.Matches()
		countFore := count(matches, fore) + 1
		countBack := count(matches, back) + 1

		p := float64(countBack) / float64(totalBack)
		return binom.ComplementCdf(countFore, totalFore, p), ""
	}
}
