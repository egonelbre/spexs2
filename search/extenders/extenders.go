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
		token := db.FullSequence[p]
		if token == ZeroToken {
			continue
		}

		q, ok := querys[token]
		if !ok {
			q = NewQuery(base, RegToken{token, IsSingle})
			querys[token] = q
		}

		q.Loc.Add(p + 1)
	}
}

type Simple struct {
	extender
}

func (e *Simple) Extend(base *Query) Querys {
	querys := make(queryMap)
	extend(base, e.Db, querys)
	return toQuerys(querys)
}

func combine(base *Query, db *Database, querys queryMap, flags RegFlags) {
	for _, group := range db.Groups {
		q := NewQuery(base, RegToken{group.Token, IsGroup | flags})
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

type Group struct {
	extender
}

func (e *Group) Extend(base *Query) Querys {
	querys := make(queryMap)
	extend(base, e.Db, querys)
	combine(base, e.Db, querys, IsSingle)
	return toQuerys(querys)
}

func starExtendPosition(base *Query, db *Database, querys queryMap, p int) {
	token := db.FullSequence[p]
	for token != ZeroToken {
		q, ok := querys[token]
		if !ok {
			q = NewQuery(base, RegToken{token, IsStar})
			querys[token] = q
		}
		q.Loc.Add(p + 1)
		token = db.FullSequence[p+1]
	}
}

func starExtend(base *Query, db *Database, querys queryMap) {
	prevseq := -1
	for _, p := range base.Loc.Iter() {
		seq := &db.PosToSequence[p]
		if seq.Index != prevseq {
			starExtendPosition(base, db, querys, p)
		}
		prevseq = seq.Index
	}
}

type Star struct {
	extender
}

func (e *Star) Extend(base *Query) Querys {
	patterns := make(queryMap)
	extend(base, e.Db, patterns)
	stars := make(queryMap)
	starExtend(base, e.Db, stars)
	return append(toQuerys(patterns), toQuerys(stars)...)
}

type Regex struct {
	extender
}

func (e *Regex) Extend(base *Query) Querys {
	patterns := make(queryMap)
	extend(base, e.Db, patterns)
	combine(base, base.Db, patterns, IsSingle)
	stars := make(queryMap)
	starExtend(base, e.Db, stars)
	combine(base, base.Db, stars, IsStar)
	return append(toQuerys(patterns), toQuerys(stars)...)
}
