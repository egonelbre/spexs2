package extenders

import "github.com/egonelbre/spexs2/search"

func starGreedyExtend(base *search.Query, db *search.Database, querys queryMap) {
	positions := base.Loc.Iter()
	if len(positions) == 0 {
		return
	}

	// initialize the last position and sequence index
	prevpos := positions[0]
	prevseq := db.PosToSequence[prevpos].Index

	for _, p := range positions {
		seq := db.PosToSequence[p]

		// if we encounter a sequence index change the prevpos was the last position in sequence
		if seq.Index == prevseq {
			prevpos = p
			continue
		}

		starExtendPosition(base, db, querys, prevpos)

		prevseq = seq.Index
		prevpos = p
	}

	// also change the last position
	starExtendPosition(base, db, querys, prevpos)
}

func StarGreedy(base *search.Query) search.Querys {
	patterns := make(queryMap)
	extend(base, base.Db, patterns)
	stars := make(queryMap)
	starGreedyExtend(base, base.Db, stars)
	return append(toQuerys(patterns), toQuerys(stars)...)
}

func RegexGreedy(base *search.Query) search.Querys {
	patterns := make(queryMap)
	extend(base, base.Db, patterns)
	combine(base, base.Db, patterns, search.IsSingle)
	stars := make(queryMap)
	starGreedyExtend(base, base.Db, stars)
	combine(base, base.Db, stars, search.IsStar)
	return append(toQuerys(patterns), toQuerys(stars)...)
}
