package extenders

import (
	. "spexs"
	"set/multi"
)

type queryMap map[Token]*Query

//TODO: method queryMap.toQuerys
//TODO: method queryMap.addLoc(token, pos)
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
	for _, p := range base.Loc.Iter() {
		token, ok, next := db.GetToken(p)
		if !ok {
			continue
		}

		q, ok := querys[token]
		if !ok {
			q = NewQuery(base, RegToken{token, false, false})
			querys[token] = q
		}

		q.Loc.Add(next)
	}
}

func Simple(base *Query) Querys {
	querys := make(queryMap)
	extend(base, base.Db, querys)
	return toQuerys(querys)
}

func combine(base *Query, db *Database, querys queryMap, isStar bool) {
	for _, group := range db.Groups {
		q := NewQuery(base, RegToken{group.Token, true, isStar})
		querys[group.Token] = q
		sets := multi.New()
		for _, token := range group.Elems {
			if s, ok := querys[token]; ok {
				sets.AddSet(s.Loc)
			}
		}
		q.Loc = sets
	}
}


func Group(base *Query) Querys {
	querys := make(queryMap)
	extend(base, base.Db, querys)
	combine(base, base.Db, querys, false)
	return toQuerys(querys)
}

func starGreedyExtend(base *Query, db *Database, querys queryMap) {
	lastPos := make(map[int]int, base.Loc.Len())
	
	for _, p := range base.Loc.Iter() {
		si := db.PosToSequence[p]
		v, ok := lastPos[si]
		if !ok || v < p {
			lastPos[si] = p
		}
	}

	for _, p := range lastPos {
		var q *Query
		token, ok, next := db.GetToken(p)
		for ok {
			q, ok = querys[token]
			if !ok {
				q = NewQuery(base, RegToken{token, false, true})
				querys[token] = q
			}
			q.Loc.Add(next)
			token, ok, next = db.GetToken(next)
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

func starExtend(base *Query, db *Database, querys queryMap) {
	firstPos := make(map[int]int, base.Loc.Len())

	for _, p := range base.Loc.Iter() {
		si := db.PosToSequence[p]
		v, ok := firstPos[si]
		if !ok || v > p {
			firstPos[si] = p
		}
	}

	for _, p := range firstPos {
		var q *Query
		token, ok, next := db.GetToken(p)
		for ok {
			q, ok = querys[token]
			if !ok {
				q = NewQuery(base, RegToken{token, false, true})
				querys[token] = q
			}
			q.Loc.Add(next)
			token, ok, next = db.GetToken(next)
		}
	}
}

func Star(base *Query) Querys {
	patterns := make(queryMap)
	extend(base, base.Db, patterns)
	stars := make(queryMap)
	starExtend(base, base.Db, stars)
	return append(toQuerys(patterns), toQuerys(stars)...)
}

func Regex(base *Query) Querys {
	patterns := make(queryMap)
	extend(base, base.Db, patterns)
	combine(base, base.Db, patterns, false)
	stars := make(queryMap)
	starExtend(base, base.Db, stars)
	combine(base, base.Db, stars, true)
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
