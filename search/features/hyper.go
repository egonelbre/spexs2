package features

import (
	. "github.com/egonelbre/spexs2/search"
	"github.com/egonelbre/spexs2/stats/hyper"
)

// hypergeometric split p-value
type Hyper struct {
	feature
	Fore []int
	Back []int
}

func (f *Hyper) Evaluate(q *Query) (float64, string) {
	totalFore := count(f.Db.Total, f.Fore)
	totalBack := count(f.Db.Total, f.Back)

	matches := q.Matches(f.Db)
	countFore := count(matches, f.Fore)
	countBack := count(matches, f.Back)

	return hyper.ComplementCdf(countFore, countBack, totalFore, totalBack), ""
}

// approximate hypergeometric split p-value (~5 significant digits)
type HyperApprox struct {
	feature
	Fore []int
	Back []int
}

func (f *HyperApprox) Evaluate(q *Query) (float64, string) {
	totalFore := count(f.Db.Total, f.Fore)
	totalBack := count(f.Db.Total, f.Back)

	matches := q.Matches(f.Db)
	countFore := count(matches, f.Fore)
	countBack := count(matches, f.Back)

	return hyper.ComplementCdfApprox(countFore, countBack, totalFore, totalBack), ""
}
