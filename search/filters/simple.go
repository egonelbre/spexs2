package filters

import (
	"encoding/json"
	"strings"

	. "github.com/egonelbre/spexs2/search"
)

// don't allow start to be grouping token
type NoStartingGroup struct {
	filter
}

func (f *NoStartingGroup) Accepts(q *Query) bool {
	start := q.Pat[0]
	return start.Flags == IsSingle

}

// don't allow ending to be grouping token
type NoEndingGroup struct {
	filter
}

func (f *NoEndingGroup) Accepts(q *Query) bool {
	end := q.Pat[len(q.Pat)-1]
	return end.Flags&IsGroup == 0

}

type NoTokens struct {
	filter
	contains map[Token]bool
}

func (f *NoTokens) Init(data []byte) {
	var filter struct{ Tokens string }
	err := json.Unmarshal(data, &filter)
	if err != nil {
		panic(err)
	}

	line := strings.TrimSpace(filter.Tokens)
	tokenNames := strings.Split(line, f.Db.Separator)
	tokens := f.Db.ToTokens(tokenNames)

	f.contains = make(map[Token]bool, len(tokens))
	for _, token := range tokens {
		f.contains[token] = true
	}
}

func (f *NoTokens) Accepts(q *Query) bool {
	e := q.Pat[len(q.Pat)-1]
	return !f.contains[e.Token]

}
