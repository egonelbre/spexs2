package features

// find optimal hypergeometric split
/*

"math"
"sort"
"stats/hyper"

func HyperOptimal(fore []int) FeatureFunc {
	return func(q *Query) (float64, string) {
		totalFore := countOnly(q.Db.Total, fore)
		totalBack := countOnly(q.Db.Total, back)

		matches := q.MatchSeqs()
		countFore := countOnly(matches, fore)
		countBack := countOnly(matches, back)

		return hyper.SplitDown(countFore, countBack, totalFore, totalBack), ""
	}
}

type optimalSplit struct {
	pvalue  float64
	matches int
	seqs    int
}



type uintSlice []uint

func (p uintSlice) Len() int           { return len(p) }
func (p uintSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p uintSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func uniq(data []uint) []uint {
	if len(data) <= 0 {
		return data
	}
	k := 1
	for i := 0; i < len(data); i += 1 {
		if data[k-1] != data[i] {
			data[k] = data[i]
			k += 1
		}
	}

	return data[0:k]
}

func (q *Query) FindOptimalSplit() float64 {
	if q.cache.optimalSplit.pvalue < 0 {
		db := q.Db
		positions := make([]uint, q.Loc.Len())
		k := 0
		for _, val := range q.Loc.Iter() {
			p, _ := DecodePos(val)
			positions[k] = p
			k += 1
		}
		sort.Sort(uintSlice(positions))
		positions = uniq(positions)

		matches := 0
		for _, c := range q.MatchSeqs() {
			matches += c
		}

		all := 0
		for _, s := range db.Sections {
			all += s.Count
		}

		accCount := 0
		splt := optimalSplit{math.Inf(1.0), -1, -1}

		for _, i := range positions {
			seq := db.Sequences[i]
			accCount += seq.Count
			p := hyper.Split(accCount, matches, int(i+1), all)
			if p < splt.pvalue {
				splt = optimalSplit{p, accCount, int(i + 1)}
			}
		}
		q.cache.optimalSplit = splt
	}
	return q.cache.optimalSplit.pvalue
}

func (q *Query) FindOptimalSplitSeqs() int {
	if q.cache.optimalSplit.pvalue < 0 {
		q.FindOptimalSplit()
	}
	return q.cache.optimalSplit.seqs
}

func (q *Query) FindOptimalSplitMatches() int {
	if q.cache.optimalSplit.pvalue < 0 {
		q.FindOptimalSplit()
	}
	return q.cache.optimalSplit.matches
}


*/
