package search

import (
	"strings"

	"github.com/egonelbre/spexs2/set"
	"github.com/egonelbre/spexs2/set/array"
	"github.com/egonelbre/spexs2/set/packed"
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
	if parent != nil {
		q.Pat = make([]RegToken, len(parent.Pat)+1)
		copy(q.Pat, parent.Pat)
		q.Pat[len(q.Pat)-1] = token
		q.Db = parent.Db
	}

	if parent == nil {
		// the initial location is only stored once,
		// so we can just use the trivial implementaton
		q.Loc = array.New()
	} else {
		// we try to guess what would be the most efficient way to represent the set
		// when the positions are sparse, packing doesn't give anything
		// when they are dense, simple method uses more memory than it needs to
		expected := parent.Loc.Len() / (len(parent.Db.Alphabet) + 1)
		spacing := len(parent.Db.PosToSequence) / (expected + 1)
		if (spacing > 1<<15) || (expected < 4) {
			q.Loc = array.New()
		} else {
			q.Loc = packed.New()
		}
	}

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
	var buf strings.Builder
	for _, tok := range q.Pat {
		if tok.Flags&IsStar != 0 {
			buf.WriteString("*")
		}
		buf.WriteRune(rune(tok.Token))
	}
	return buf.String()
}

func (q *Query) string(short bool) string {
	var buf strings.Builder
	db := q.Db
	for i, tok := range q.Pat {
		if tok.Flags&IsStar != 0 {
			buf.WriteString("*")
			buf.WriteString(db.Separator)
		}

		tokInfo, ok := db.Alphabet[tok.Token]
		if ok {
			buf.WriteString(tokInfo.Name)
		}

		group, ok := db.Groups[tok.Token]
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

	return buf.String()
}
