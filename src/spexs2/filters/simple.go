package filters

import (
	"encoding/json"
	. "spexs"
	"strings"
)

// don't allow start or end to be grouping token
func NoGroupingEnds(s Setup, data []byte) FilterFunc {
	return func(q *Query) bool {
		start := q.Pat[0]
		end := q.Pat[len(q.Pat)-1]
		return !start.IsGroup && !start.IsStar && !end.IsGroup
	}
}

func NoTokens(s Setup, data []byte) FilterFunc {
	var filter struct{ Tokens string }
	err := json.Unmarshal(data, &filter)
	if err != nil {
		panic(err)
	}

	line := strings.TrimSpace(filter.Tokens)
	tokenNames := strings.Split(line, setup.Db.Separator)
	tokens := setup.Db.ToTokens(tokenNames)

	contains := make(map[Token]bool, len(tokens))
	for _, token := range tokens {
		contains[token] = true
	}

	return func(q *Query) bool {
		e := q.Pat[len(q.Pat)-1]
		return !contains[e.Token]
	}
}
