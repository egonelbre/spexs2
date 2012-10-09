package spexs

import (
	"bytes"
	"math"
	"sort"
	"spexs/sets"
	"stats/hyper"
)

type RegToken struct {
	Token   Token
	IsGroup bool
	IsStar  bool
}

type Query struct {
	Pat []RegToken
	Loc *sets.HashSet

	cache queryCache
}

func EncodePos(idx int, pos int) int {
	return (idx << 8) | (pos & 0xFF)
}

func DecodePos(val int) (int, int) {
	return val >> 8, val & 0xFF
}

func NewQuery(parent *Query, token RegToken) *Query {
	q := &Query{}

	if parent != nil {
		size := len(parent.Pat) + 1
		q.Pat = make([]RegToken, size)
		copy(q.Pat, parent.Pat)
		q.Pat[size-1] = token

		estimatedSize := parent.Loc.Len() / 8
		q.Loc = sets.NewHashSet(estimatedSize)
	} else {
		q.Pat = nil
		q.Loc = sets.NewHashSet(0)
	}

	q.cache.reset()

	return q
}

func NewEmptyQuery(db *Database) *Query {
	q := NewQuery(nil, RegToken{})
	for i, _ := range db.Sequences {
		last := 0
		_, ok, next := db.GetToken(i, last)
		for ok {
			q.Loc.Add(EncodePos(i, last))
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
	count        []int
	occs         []int
	optimalSplit optimalSplit
}

type optimalSplit struct {
	pvalue  float64
	matches int
	seqs    int
}

func (q *queryCache) reset() {
	q.count = nil
	q.occs = nil
	q.optimalSplit.pvalue = -1.0
}

func (q *Query) CacheValues(db *Database) {
	if q.cache.count == nil {
		q.SeqCount(db)
	}
	if q.cache.occs == nil {
		q.MatchCount(db)
	}
	q.Loc.Clear()
}

func (q *Query) SeqCount(db *Database) []int {
	if q.cache.count == nil {
		count := make([]int, len(db.Sections))

		for val := range q.Loc.Iter() {
			i, _ := DecodePos(val)
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

		for val := range q.Loc.Iter() {
			i, _ := DecodePos(val)
			seq := db.Sequences[i]
			occs[seq.Section] += seq.Count
		}

		q.cache.occs = occs
	}
	return q.cache.occs
}

func (q *Query) FindOptimalSplit(db *Database) float64 {
	if q.cache.optimalSplit.pvalue < 0 {
		positions := make([]int, q.Loc.Len())
		k := 0
		for i := range q.Loc.Iter() {
			positions[k] = i
			k += 1
		}
		sort.Ints(positions)

		matches := 0
		for _, c := range q.SeqCount(db) {
			matches += c
		}

		all := 0
		for _, s := range db.Sections {
			all += s.Count
		}

		accCount := 0
		splt := optimalSplit{math.Inf(1.0), -1, -1}

		for _, i := range positions {
			seq := db.Sequences[i]
			accCount += seq.Count
			p := hyper.Split(accCount, matches, i+1, all)
			if p < splt.pvalue {
				splt = optimalSplit{p, accCount, i + 1}
			}
		}
		q.cache.optimalSplit = splt
	}
	return q.cache.optimalSplit.pvalue
}

func (q *Query) FindOptimalSplitSeqs(db *Database) int {
	if q.cache.optimalSplit.pvalue < 0 {
		q.FindOptimalSplit(db)
	}
	return q.cache.optimalSplit.seqs
}

func (q *Query) FindOptimalSplitMatches(db *Database) int {
	if q.cache.optimalSplit.pvalue < 0 {
		q.FindOptimalSplit(db)
	}
	return q.cache.optimalSplit.matches
}

func (q *Query) String(db *Database, short bool) string {
	buf := bytes.NewBufferString("")

	for i, regToken := range q.Pat {

		if regToken.IsStar {
			buf.WriteString("*")
			buf.WriteString(db.Separator)
		}

		tokInfo, ok := db.Alphabet[regToken.Token]
		if ok {
			buf.WriteString(tokInfo.Name)
		}

		group, ok := db.Groups[regToken.Token]
		if ok {
			if short {
				buf.WriteString(group.Name)
			} else {
				buf.WriteString(group.FullName)
			}
		}

		isLast := len(q.Pat)-1 == i
		if !isLast {
			buf.WriteString(db.Separator)
		}
	}

	return string(buf.Bytes())
}
