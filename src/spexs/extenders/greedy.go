package extenders

import . "spexs"

func starGreedyExtend(base *Query, db *Database, querys queryMap) {
	if base.Loc.IsSorted() {
		positions := base.Loc.Iter()
		if len(positions) == 0 {
			return
		}

		// initialize the last position and sequence index
		last_p := positions[0]
		last_si := db.PosToSequence[last_p].Index

		for _, p := range positions {
			seq := db.PosToSequence[p]

			// if we encounter a sequence index change the last_p was the last position in sequence
			if seq.Index == last_si {
				last_p = p
				continue
			}

			starExtendPosition(base, db, querys, last_p)

			last_si = seq.Index
			last_p = p
		}

		// also change the last position
		starExtendPosition(base, db, querys, last_p)

	} else {
		lastPos := make(map[int]int, base.Loc.Len())
		for _, p := range base.Loc.Iter() {
			seq := db.PosToSequence[p]
			v, ok := lastPos[seq.Index]
			if !ok || v < p {
				lastPos[seq.Index] = p
			}
		}

		for _, p := range lastPos {
			starExtendPosition(base, db, querys, p)
		}
	}
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
	combine(base, base.Db, patterns, false)
	stars := make(queryMap)
	starGreedyExtend(base, base.Db, stars)
	combine(base, base.Db, stars, true)
	return append(toQuerys(patterns), toQuerys(stars)...)
}
