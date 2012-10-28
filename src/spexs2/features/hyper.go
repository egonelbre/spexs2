package features

import (
	. "spexs"
	"stats/hyper"
)

// hypergeometric split p-value
func Hyper(fore []int, back []int) FeatureFunc {
	return func(q *Query) (float64, string) {
		totalFore := countOnly(q.Db.Total, fore)
		totalBack := countOnly(q.Db.Total, back)

		matches := q.Matches()
		countFore := countOnly(matches, fore)
		countBack := countOnly(matches, back)

		return hyper.Split(countFore, countBack, totalFore, totalBack), ""
	}
}

// approximate hypergeometric split p-value (~5 significant digits)
func HyperApprox(fore []int, back []int) FeatureFunc {
	return func(q *Query) (float64, string) {
		totalFore := countOnly(q.Db.Total, fore)
		totalBack := countOnly(q.Db.Total, back)

		matches := q.Matches()
		countFore := countOnly(matches, fore)
		countBack := countOnly(matches, back)

		return hyper.SplitApprox(countFore, countBack, totalFore, totalBack), ""
	}
}

// hypergeometric split down p-value
func HyperDown(fore []int, back []int) FeatureFunc {
	return func(q *Query) (float64, string) {
		totalFore := countOnly(q.Db.Total, fore)
		totalBack := countOnly(q.Db.Total, back)

		matches := q.Matches()
		countFore := countOnly(matches, fore)
		countBack := countOnly(matches, back)

		return hyper.SplitDown(countFore, countBack, totalFore, totalBack), ""
	}
}
