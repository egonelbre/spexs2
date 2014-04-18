package features

import . "github.com/egonelbre/spexs2/search"

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
type Total struct {
	feature
	Group []int
}

func (f *Total) Evaluate(q *Query) (float64, string) {
	return countf(f.Db.Total, f.Group), ""
}

// the count of matching sequences
type Matches struct {
	feature
	Group []int
}

func (f *Matches) Evaluate(q *Query) (float64, string) {
	return countf(q.Matches(f.Db), f.Group), ""
}

const minSeqsCountTable = 30

// the count of matching unique sequences
type Seqs struct {
	feature
	Group []int
}

func (f *Seqs) Evaluate(q *Query) (float64, string) {
	prevseq := -1

	count := make([]int, len(f.Db.Total))
	for _, p := range q.Loc.Iter() {
		seq := f.Db.PosToSequence[p]
		if seq.Index == prevseq {
			continue
		}
		prevseq = seq.Index
		count[seq.Section] += 1
	}
	return countf(count, f.Group), ""
}

// the count of occurences
type Occs struct {
	feature
	Group []int
}

func (f *Occs) Evaluate(q *Query) (float64, string) {
	occs := q.Occs(f.Db)
	return countf(occs, f.Group), ""
}

// the ratio of matching sequences to total count
type MatchesProp struct {
	feature
	Group []int
}

func (f *MatchesProp) Evaluate(q *Query) (float64, string) {
	matches := q.Matches(f.Db)
	total := countf(f.Db.Total, f.Group)
	return countf(matches, f.Group) / total, ""
}

// the ratio between matching sequences (adjusted)
type MatchesRatio struct {
	feature
	Nom   []int
	Denom []int
}

func (f *MatchesRatio) Evaluate(q *Query) (float64, string) {
	matches := q.Matches(f.Db)
	countNom := countf(matches, f.Nom) + 1.0
	countDenom := countf(matches, f.Denom) + 1.0
	return countNom / countDenom, ""
}

// the ratio between occurences (adjusted)
type OccsRatio struct {
	feature
	Nom   []int
	Denom []int
}

func (f *OccsRatio) Evaluate(q *Query) (float64, string) {
	occs := q.Occs(f.Db)
	countNom := countf(occs, f.Nom) + 1.0
	countDenom := countf(occs, f.Denom) + 1.0
	return countNom / countDenom, ""
}

// the ratio of proptions between matches (adjusted)
type MatchesPropRatio struct {
	feature
	Nom   []int
	Denom []int
}

func (f *MatchesPropRatio) Evaluate(q *Query) (float64, string) {
	totalNom := countf(f.Db.Total, f.Nom) + 1.0
	totalDenom := countf(f.Db.Total, f.Denom) + 1.0

	matches := q.Matches(f.Db)
	countNom := countf(matches, f.Nom) + 1.0
	countDenom := countf(matches, f.Denom) + 1.0
	return (countNom / totalNom) / (countDenom / totalDenom), ""

}
