package extenders

import (
	. "github.com/egonelbre/spexs2/search"
	"github.com/egonelbre/spexs2/set/multi"
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

func starExtendPosition(base *Query, db *Database, querys queryMap, p int) {
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

func starExtend(base *Query, db *Database, querys queryMap) {
	prevseq := -1
	for _, p := range base.Loc.Iter() {
		seq := db.PosToSequence[p]
		if seq.Index == prevseq {
			continue
		}
		starExtendPosition(base, db, querys, p)
		prevseq = seq.Index
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
