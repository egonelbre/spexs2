package spexs

import (
	"bytes"
	"math"
	set "set/trie"
	"sort"
	"stats/hyper"
)

type RegToken struct {
	Token   Token
	IsGroup bool
	IsStar  bool
}

type featureResult struct {
	Value float64
	Info  string
}

type Query struct {
	Pat   []RegToken
	Loc   *set.Set
	Db    *Database
	memo  map[FeatureFunc]featureResult
	cache countCache
}

var PosOffset uint = 8

func EncodePos(idx uint, pos uint) uint {
	return (idx << PosOffset) | pos
}

func DecodePos(val uint) (uint, uint) {
	return val >> PosOffset, val & ((1 << PosOffset) - 1)
}

func NewQuery(parent *Query, token RegToken) *Query {
	q := &Query{}

	q.Pat = nil
	q.Db = nil
	q.cache = make(map[FeatureFunc]featureResult)
	if parent != nil {
		q.Pat = make([]RegToken, len(parent.Pat)+1)
		copy(q.Pat, parent.Pat)
		q.Pat[len(q.Pat)-1] = token
		q.Db = parent.Db
	}
	q.Loc = set.New()
	q.cache.reset()

	return q
}

func NewEmptyQuery(db *Database) *Query {
	q := NewQuery(nil, RegToken{})
	q.Db = db
	for idx, _ := range db.Sequences {
		i := uint(idx)
		last := uint(0)
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

func (q *Query) Memoized(f FeatureFunc) (float64, string) {
	if res, ok := q.memo[f]; ok {
		return res.Value, res.Info
	}
	val, info := f(q)
	q.memo[f] = Result{val, info}
	return val, info
}

type queryCache struct {
	count []int
	occs  []int
}

func (q *countCache) reset() {
	q.count = nil
	q.occs = nil
}

func (q *Query) CacheValues() {
	q.MatchSeqs()
	q.MatchOccs()
	q.Loc = nil
}

func (q *Query) MatchSeqs() []int {
	if q.cache.count == nil {
		db := q.Db
		counted := make(map[uint]bool, q.Loc.Len())
		count := make([]int, len(db.Sections))

		for _, val := range q.Loc.Iter() {
			i, _ := DecodePos(val)
			if counted[i] {
				continue
			}
			counted[i] = true
			seq := db.Sequences[i]
			count[seq.Section] += seq.Count
		}

		q.cache.count = count
	}
	return q.cache.count
}

func (q *Query) MatchOccs() []int {
	if q.cache.occs == nil {
		db := q.Db
		occs := make([]int, len(db.Sections))

		for _, val := range q.Loc.Iter() {
			i, _ := DecodePos(val)
			seq := db.Sequences[i]
			occs[seq.Section] += seq.Count
		}

		q.cache.occs = occs
	}
	return q.cache.occs
}

func (q *Query) String() string {
	return q.string(true)
}

func (q *Query) StringLong() string {
	return q.string(false)
}

func (q *Query) string(short bool) string {
	buf := bytes.NewBufferString("")
	db := q.Db
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
