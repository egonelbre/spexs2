package spexs

import (
	"bytes"
	"spexs/sets"
	"utils"
)

type Rid struct {
	Id      Tid
	IsGroup bool
	IsStar  bool
}

type Query struct {
	Pat []Rid
	Loc *sets.HashSet

	cache queryCache
}

func NewQuery(parent *Query, token Rid) *Query {
	q := &Query{}

	if parent != nil {
		q.Pat = append(parent.Pat, token)

		estimatedSize := parent.Loc.Len() / 8
		q.Loc = sets.NewHashSet(estimatedSize)
	} else {
		q.Pat = nil
		q.Loc = sets.NewHashSet(0)
	}

	return q
}

func NewEmptyQuery(db *Database) *Query {
	q := NewQuery(nil, Rid{})
	for i, _ := range db.Sequences {
		last := 0
		_, ok, next := db.GetToken(i, last)
		for ok {
			q.Loc.Add(i, last)
			last = next
			_, ok, next = db.GetToken(i, next)
		}
	}
	return q
}

func (q *Query) Len() int {
	return len(q.Pat)
}

type queryCache struct {
	count []int
	occs  []int
}

func (q *Query) SeqCount(db *Database) []int {
	if q.cache.count == nil {
		count := make([]int, len(db.Sections))

		for i := range q.Loc.Iter() {
			seq := db.Sequences[i]
			count[seq.Section] += seq.Count
		}

		q.cache.count = count
	}
	return q.cache.count
}

func (q *Query) MatchCount(db *Database) []int {
	if q.cache.occs == nil {
		occs := make([]int, len(db.Sections))

		for i, pv := range q.Loc.Iter() {
			seq := db.Sequences[i]
			matchCount := utils.BitCount64(uint64(pv))
			occs[seq.Section] += seq.Count * matchCount
		}

		q.cache.occs = occs
	}
	return q.cache.occs
}

func (q *Query) String(db *Database, short bool) string {
	buf := bytes.NewBufferString("")

	for _, rid := range q.Pat {
		alpha, ok := db.Alphabet[rid.Id]
		if ok {
			buf.WriteString(alpha.Str)
		}

		group, ok := db.Groups[rid.Id]
		if ok {
			if short {
				buf.WriteString(group.Alias)
			} else {
				buf.WriteString(group.Str)
			}

		}
	}

	return string(buf.Bytes())
}
