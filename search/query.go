package search

import (
	"bytes"
	"github.com/egonelbre/spexs2/set"
	defset "github.com/egonelbre/spexs2/set/rle"
)

type RegFlags uint8

const (
	IsGroup RegFlags = 1 << iota
	IsStar
	IsSingle RegFlags = 0
)

type RegToken struct {
	Token Token
	Flags RegFlags
}

type Query struct {
	Pat   []RegToken
	Loc   set.Set
	Db    *Database
	cache queryCache
}

func NewQuery(parent *Query, token RegToken) *Query {
	q := &Query{}

	q.Pat = nil
	if parent != nil {
		q.Pat = make([]RegToken, len(parent.Pat)+1)
		copy(q.Pat, parent.Pat)
		q.Pat[len(q.Pat)-1] = token
	}
	q.Loc = defset.New()
	q.cache.reset()

	return q
}

func NewEmptyQuery(db *Database) *Query {
	q := NewQuery(nil, RegToken{})
	db.AddAllPositions(q.Loc)
	return q
}

func (q *Query) Len() int {
	return len(q.Pat)
}

type queryCache struct {
	count []int
	occs  []int
}

func (q *queryCache) reset() {
	q.count = nil
	q.occs = nil
}

func (q *Query) CacheValues(db *Database) {
	if q.cache.count == nil || q.cache.occs == nil {
		q.cache.count, q.cache.occs = db.MatchesOccs(q.Loc)
	}
}

func (q *Query) Matches(db *Database) []int {
	q.CacheValues(db)
	return q.cache.count
}

func (q *Query) Occs(db *Database) []int {
	q.CacheValues(db)
	return q.cache.occs
}

func (q *Query) String(db *Database) string {
	return q.string(db, true)
}

func (q *Query) StringLong(db *Database) string {
	return q.string(db, false)
}

func (q *Query) StringRaw() string {
	buf := bytes.NewBufferString("")
	for _, tok := range q.Pat {
		buf.WriteRune(rune(tok.Token))
	}
	return string(buf.Bytes())
}

func (q *Query) string(db *Database, short bool) string {
	buf := bytes.NewBufferString("")
	for i, regToken := range q.Pat {
		if regToken.Flags&IsStar != 0 {
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
