package filters

import (
	"encoding/json"
	. "github.com/egonelbre/spexs2/search"
	"strings"
)

// don't allow start to be grouping token
func NoStartingGroup(s Setup, data []byte) Filter {
	return func(q *Query) bool {
		start := q.Pat[0]
		return !start.IsGroup && !start.IsStar
	}
}

// don't allow ending to be grouping token
func NoEndingGroup(s Setup, data []byte) Filter {
	return func(q *Query) bool {
		end := q.Pat[len(q.Pat)-1]
		return !end.IsGroup
	}
}

func NoTokens(s Setup, data []byte) Filter {
	var filter struct{ Tokens string }
	err := json.Unmarshal(data, &filter)
	if err != nil {
		panic(err)
	}

	line := strings.TrimSpace(filter.Tokens)
	tokenNames := strings.Split(line, s.Db.Separator)
	tokens := s.Db.ToTokens(tokenNames)

	contains := make(map[Token]bool, len(tokens))
	for _, token := range tokens {
		contains[token] = true
	}

	return func(q *Query) bool {
		e := q.Pat[len(q.Pat)-1]
		return !contains[e.Token]
	}
}
