package features

import . "spexs"

// function to sum elements in arr by group
func count(arr []int, group []int) int {
	total := 0
	for _, id := range group {
		total += arr[id]
	}
	return total
}

func countf(arr []int, group []int) float64 {
	return float64(count(arr, group))
}

// the total count of sequences
func Total(group []int) Feature {
	return func(q *Query) (float64, string) {
		total := countf(q.Db.Total, group)
		return total, ""
	}
}

// the count of matching sequences
func Matches(group []int) Feature {
	return func(q *Query) (float64, string) {
		matches := q.Matches()
		return countf(matches, group), ""
	}
}

// the count of matching unique sequences
func Seqs(group []int) Feature {
	return func(q *Query) (float64, string) {
		db := q.Db
		counted := make(map[uint]bool, 30)
		count := 0
		for _, val := range q.Loc.Iter() {
			i, _ := DecodePos(val)
			if counted[i] {
				continue
			}
			counted[i] = true
			seq := db.Sequences[i]
			for _, sec := range group {
				if _, ok := seq.Count[sec]; ok {
					count += 1
					break
				}
			}
		}
		return float64(count), ""
	}
}

// the count of occurences
func Occs(group []int) Feature {
	return func(q *Query) (float64, string) {
		occs := q.Occs()
		return countf(occs, group), ""
	}
}

// the ratio of matching sequences to total count
func MatchesProp(group []int) Feature {
	return func(q *Query) (float64, string) {
		total := countf(q.Db.Total, group)
		matches := q.Matches()
		return countf(matches, group) / total, ""
	}
}

// the ratio between matching sequences (adjusted)
func MatchesRatio(nom []int, denom []int) Feature {
	return func(q *Query) (float64, string) {
		matches := q.Matches()
		countNom := countf(matches, nom) + 1.0
		countDenom := countf(matches, denom) + 1.0
		return countNom / countDenom, ""
	}
}

// the ratio between occurences (adjusted)
func OccsRatio(nom []int, denom []int) Feature {
	return func(q *Query) (float64, string) {
		occs := q.Occs()
		countNom := countf(occs, nom) + 1.0
		countDenom := countf(occs, denom) + 1.0
		return countNom / countDenom, ""
	}
}

// the ratio of proptions between matches (adjusted)
func MatchesPropRatio(nom []int, denom []int) Feature {
	return func(q *Query) (float64, string) {
		totalNom := countf(q.Db.Total, nom) + 1.0
		totalDenom := countf(q.Db.Total, denom) + 1.0

		matches := q.Matches()
		countNom := countf(matches, nom) + 1.0
		countDenom := countf(matches, denom) + 1.0
		return (countNom / totalNom) / (countDenom / totalDenom), ""
	}
}
