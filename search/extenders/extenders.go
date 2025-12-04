package extenders

import (
	"github.com/egonelbre/spexs2/search"
	"github.com/egonelbre/spexs2/set/multi"
)

type queryMap map[search.Token]*search.Query

// TODO: method queryMap.toQuerys
// TODO: method queryMap.addLoc(token, pos)
func toQuerys(queryMap queryMap) search.Querys {
	querys := make(search.Querys, len(queryMap))
	i := 0
	for _, q := range queryMap {
		querys[i] = q
		i++
	}
	return querys
}

func extend(base *search.Query, db *search.Database, querys queryMap) {
	for _, p := range base.Loc.All() {
		token := db.FullSequence[p]
		if token == search.ZeroToken {
			continue
		}

		q, ok := querys[token]
		if !ok {
			q = search.NewQuery(base, search.RegToken{Token: token, Flags: search.IsSingle})
			querys[token] = q
		}

		q.Loc.Add(p + 1)
	}
}

func Simple(base *search.Query) search.Querys {
	querys := make(queryMap)
	extend(base, base.Db, querys)
	return toQuerys(querys)
}

func combine(base *search.Query, db *search.Database, querys queryMap, flags search.RegFlags) {
	for _, group := range db.Groups {
		q := search.NewQuery(base, search.RegToken{Token: group.Token, Flags: search.IsGroup | flags})
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

func Group(base *search.Query) search.Querys {
	querys := make(queryMap)
	extend(base, base.Db, querys)
	combine(base, base.Db, querys, search.IsSingle)
	return toQuerys(querys)
}

func starExtendPosition(base *search.Query, db *search.Database, querys queryMap, start int) {
	for i, token := range db.FullSequence[start:] {
		p := start + i
		if token == search.ZeroToken {
			break
		}

		q, ok := querys[token]
		if !ok {
			q = search.NewQuery(base, search.RegToken{Token: token, Flags: search.IsStar})
			querys[token] = q
		}

		q.Loc.Add(p + 1)
	}
}

func starExtend(base *search.Query, db *search.Database, querys queryMap) {
	prevseq := -1
	for _, p := range base.Loc.All() {
		seq := db.PosToSequence[p]
		if seq.Index == prevseq {
			continue
		}
		starExtendPosition(base, db, querys, p)
		prevseq = seq.Index
	}
}

func Star(base *search.Query) search.Querys {
	patterns := make(queryMap)
	extend(base, base.Db, patterns)
	stars := make(queryMap)
	starExtend(base, base.Db, stars)
	return append(toQuerys(patterns), toQuerys(stars)...)
}

func Regex(base *search.Query) search.Querys {
	patterns := make(queryMap)
	extend(base, base.Db, patterns)
	combine(base, base.Db, patterns, search.IsSingle)
	stars := make(queryMap)
	starExtend(base, base.Db, stars)
	combine(base, base.Db, stars, search.IsStar)
	return append(toQuerys(patterns), toQuerys(stars)...)
}
