package spexs

import (
	"bytes"
	"set"
	defset "set/rle"
)

type RegToken struct {
	Token   Token
	IsGroup bool
	IsStar  bool
}

type Query struct {
	Pat   []RegToken
	Loc   set.Set
	Db    *Database
	cache queryCache
}

func NewQuery(parent *Query, token RegToken) *Query {
	q := &Query{}
	if parent != nil {
		q.Pat = make([]RegToken, len(parent.Pat)+1)
		copy(q.Pat, parent.Pat)
		q.Pat[len(q.Pat)-1] = token
		q.Db = parent.Db
	}
	q.Loc = defset.New()
	return q
}

func NewEmptyQuery(db *Database) *Query {
	q := NewQuery(nil, RegToken{})
	q.Db = db
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

func (q *Query) CacheValues() {
	if q.cache.count == nil || q.cache.occs == nil {
		q.cache.count, q.cache.occs = q.Db.MatchesOccs(q.Loc)
	}
}

func (q *Query) Matches() []int {
	q.CacheValues()
	return q.cache.count
}

func (q *Query) Occs() []int {
	q.CacheValues()
	return q.cache.occs
}

func (q *Query) String() string {
	return q.string(true)
}

func (q *Query) StringLong() string {
	return q.string(false)
}

func (q *Query) StringRaw() string {
	buf := bytes.NewBufferString("")
	for _, tok := range q.Pat {
		buf.WriteRune(rune(tok.Token))
	}
	return string(buf.Bytes())
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
