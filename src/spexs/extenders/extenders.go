package extenders

import (
	. "spexs"
	"utils"
)

type queryMap map[Token]*Query

func toQuerys(queryMap queryMap) Querys {
	querys := make(Querys, len(queryMap))
	i := 0
	for _, q := range queryMap {
		querys[i] = q
		i += 1
	}
	return querys
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
				q = NewQuery(base, RegToken{token, false, false})
				querys[token] = q
			}
			q.Loc.Add(i, next)
		}
	}
}

func Simplex(base *Query, db *Database) Querys {
	querys := make(queryMap)
	extend(base, db, querys)
	return toQuerys(querys)
}

func combine(base *Query, db *Database, querys queryMap, isStar bool) {
	for _, group := range db.Groups {
		q := NewQuery(base, RegToken{group.Token, true, isStar})
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
	return toQuerys(querys)
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
				q = NewQuery(base, RegToken{token, false, true})
				querys[token] = q
			}
			q.Loc.Add(i, next)
			token, ok, next = db.GetToken(i, next)
		}
	}
}

func Starex(base *Query, db *Database) Querys {
	patterns := make(queryMap)
	extend(base, db, patterns)
	stars := make(queryMap)
	starExtend(base, db, stars)
	return append(toQuerys(patterns), toQuerys(stars)...)
}

func Regex(base *Query, db *Database) Querys {
	patterns := make(queryMap)
	extend(base, db, patterns)
	combine(base, db, patterns, false)
	stars := make(queryMap)
	starExtend(base, db, stars)
	combine(base, db, stars, true)
	return append(toQuerys(patterns), toQuerys(stars)...)
}
