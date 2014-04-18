package extenders

import . "github.com/egonelbre/spexs2/search"

func starGreedyExtend(base *Query, db *Database, querys queryMap) {
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

func StarGreedy(base *Query) Querys {
	patterns := make(queryMap)
	extend(base, base.Db, patterns)
	stars := make(queryMap)
	starGreedyExtend(base, base.Db, stars)
	return append(toQuerys(patterns), toQuerys(stars)...)
}

func RegexGreedy(base *Query) Querys {
	patterns := make(queryMap)
	extend(base, base.Db, patterns)
	combine(base, base.Db, patterns, IsSingle)
	stars := make(queryMap)
	starGreedyExtend(base, base.Db, stars)
	combine(base, base.Db, stars, IsStar)
	return append(toQuerys(patterns), toQuerys(stars)...)
}
