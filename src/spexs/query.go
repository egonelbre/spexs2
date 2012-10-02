package spexs

import (
	"bytes"
	"fmt"
	"math"
	"sort"
	"spexs/sets"
	"stats/hyper"
	"utils"
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

	q.cache.acc = -1.0

	return q
}

func NewEmptyQuery(db *Database) *Query {
	q := NewQuery(nil, RegToken{})
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
	count    []int
	occs     []int
	acc      float64
	accdebug string
}

func (q *Query) CacheValues(db *Database) {
	if q.cache.count == nil {
		q.SeqCount(db)
	}
	if q.cache.occs == nil {
		q.MatchCount(db)
	}
	if math.IsNaN(q.cache.acc) {
		q.AccumulativeSplit(db)
	}
	q.Loc.Clear()
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

func (q *Query) AccumulativeSplit(db *Database) float64 {
	if q.cache.acc < 0 {
		matches := q.SeqCount(db)[0]
		count := db.Sections[0].Count
		acc := 0
		min := float64(10.0)
		debug := ""

		positions := make([]int, q.Loc.Len())
		k := 0
		for i := range q.Loc.Iter() {
			positions[k] = i
			k += 1
		}

		sort.Ints(positions)
		for _, i := range positions {
			seq := db.Sequences[i]
			acc += seq.Count
			p := hyper.Split(acc, matches, i+1, count)
			if p < min {
				min = p
				debug = fmt.Sprintf("hyp(%v,%v,%v,%v)", acc, matches, i+1, count)
			}
		}
		q.cache.acc = min
		q.cache.accdebug = debug
	}
	return q.cache.acc
}

func (q *Query) AccumulativeDebug(db *Database) string {
	if q.cache.acc < 0 {
		q.AccumulativeSplit(db)
	}
	return q.cache.accdebug
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
