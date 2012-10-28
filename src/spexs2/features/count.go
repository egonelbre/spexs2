package features

import . "spexs"

// the count of sequences
func Seqs(group []int) FeatureFunc {
	return func(q *Query) (float64, string) {
		total := countOnly(q.Db.Total, group)
		return total, ""
	}
}

// the count of matching sequences
func Matches(group []int) FeatureFunc {
	return func(q *Query) (float64, string) {
		matches := q.Matches()
		return countOnly(matches, group), ""
	}
}

// the count of occurences
func Occs(group []int) FeatureFunc {
	return func(q *Query) (float64, string) {
		occs := q.Occs()
		return countOnly(occs, group), ""
	}
}

// the ratio of matching sequences to total count
func MatchesProp(group []int) FeatureFunc {
	return func(q *Query) (float64, string) {
		total := countOnly(q.Db.Total, group)
		matches := q.Matches()
		return countOnly(matches, group) / total, ""
	}
}

// the ratio between matching sequences (adjusted)
func MatchesRatio(nom []int, denom []int) FeatureFunc {
	return func(q *Query) (float64, string) {
		matches := q.Matches()
		countNom := countOnly(matches, nom) + 1.0
		countDenom := countOnly(matches, denom) + 1.0
		return countNom / countDenom, ""
	}
}

// the ratio between occurences (adjusted)
func OccsRatio(nom []int, denom []int) FeatureFunc {
	return func(q *Query) (float64, string) {
		occs := q.Occs()
		countNom := countOnly(occs, nom) + 1.0
		countDenom := countOnly(occs, denom) + 1.0
		return countNom / countDenom, ""
	}
}

// the ratio of proptions between matches (adjusted)
func MatchesPropRatio(nom []int, denom []int) FeatureFunc {
	return func(q *Query) (float64, string) {
		totalNom := countOnly(q.Db.Total, nom) + 1.0
		totalDenom := countOnly(q.Db.Total, denom) + 1.0

		matches := q.Matches()
		countNom := countOnly(matches, nom) + 1.0
		countDenom := countOnly(matches, denom) + 1.0
		return (countNom / totalNom) / (countDenom / totalDenom), ""
	}
}
