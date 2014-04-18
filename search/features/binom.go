package features

import (
	. "github.com/egonelbre/spexs2/search"
	"github.com/egonelbre/spexs2/stats/binom"
)

// binomial p-value
type Binom struct {
	feature
	Fore []int
	Back []int
}

func (f *Binom) Evaluate(q *Query) (float64, string) {
	totalFore := count(f.Db.Total, f.Fore) + 1
	totalBack := count(f.Db.Total, f.Back) + 1

	matches := q.Matches(f.Db)
	countFore := count(matches, f.Fore) + 1
	countBack := count(matches, f.Back) + 1

	p := float64(countBack) / float64(totalBack)
	return binom.ComplementCdf(countFore, totalFore, p), ""
}
