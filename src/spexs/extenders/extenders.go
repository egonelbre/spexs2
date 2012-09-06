package extenders

import (
	. "spexs"
	"utils"
)

type queryMap map[Token]*Query

func flush(querys queryMap, to Querys) {
	for _, q := range querys {
		to <- q
	}
}

func extend(base *Query, db *Database, querys queryMap) {
	for i, pv := range base.Loc.Iter() {
		seq := db.Sequences[i]
		seqLen := seq.Len

		for k := 0; k < seqLen; k += 1 {
			if ((pv >> uint(k)) & 1) == 0 {
				continue
			}

			token, ok, next := db.GetToken(i, k)
			if !ok {
				continue
			}

			q, ok := querys[token]
			if !ok {
				q = NewQuery(base, Rid{token, false, false})
				querys[token] = q
			}
			q.Loc.Add(i, next)
		}
	}
}

func Simplex(base *Query, db *Database) Querys {
	querys := make(queryMap)

	extend(base, db, querys)

	result := NewQuerys()
	flush(querys, result)
	close(result)
	return result
}

func combine(base *Query, db *Database, querys queryMap, isStar bool) {
	for _, group := range db.Groups {
		q := NewQuery(base, Rid{group.Token, true, isStar})
		querys[group.Token] = q
		for _, token := range group.Elems {
			single, ok := querys[token]
			if ok {
				q.Loc.AddSet(single.Loc)
			}
		}
	}
}

func Groupex(base *Query, db *Database) Querys {
	querys := make(queryMap)

	extend(base, db, querys)
	combine(base, db, querys, false)

	result := NewQuerys()
	flush(querys, result)
	close(result)
	return result
}

func starExtend(base *Query, db *Database, querys queryMap) {
	for i, pv := range base.Loc.Iter() {
		last := utils.BitScanLeft64(uint64(pv))
		if last < 0 {
			continue
		}

		token, ok, next := db.GetToken(i, last)
		for ok {
			q, ok := querys[token]
			if !ok {
				q = NewQuery(base, Rid{token, false, true})
				querys[token] = q
			}
			q.Loc.Add(i, next)
			token, ok, next = db.GetToken(i, next)
		}
	}
}

func Starex(base *Query, db *Database) Querys {
	result := NewQuerys()

	patterns := make(queryMap)
	extend(base, db, patterns)
	flush(patterns, result)

	stars := make(queryMap)
	starExtend(base, db, stars)
	flush(stars, result)

	close(result)
	return result
}

func Regex(base *Query, db *Database) Querys {
	result := NewQuerys()

	patterns := make(queryMap)
	extend(base, db, patterns)
	combine(base, db, patterns, false)
	flush(patterns, result)

	stars := make(queryMap)
	starExtend(base, db, stars)
	combine(base, db, stars, true)
	flush(stars, result)

	close(result)
	return result
}
