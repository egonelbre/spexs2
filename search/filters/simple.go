package filters

import (
	"encoding/json"
	"strings"

	"github.com/egonelbre/spexs2/search"
)

// don't allow start to be grouping token
func NoStartingGroup(s *search.Setup, data []byte) search.Filter {
	return func(q *search.Query) bool {
		start := q.Pat[0]
		return start.Flags == search.IsSingle
	}
}

// don't allow ending to be grouping token
func NoEndingGroup(s *search.Setup, data []byte) search.Filter {
	return func(q *search.Query) bool {
		end := q.Pat[len(q.Pat)-1]
		return end.Flags&search.IsGroup == 0
	}
}

func NoTokens(s *search.Setup, data []byte) search.Filter {
	var filter struct{ Tokens string }
	err := json.Unmarshal(data, &filter)
	if err != nil {
		panic(err)
	}

	line := strings.TrimSpace(filter.Tokens)
	tokenNames := strings.Split(line, s.Db.Separator)
	tokens := s.Db.ToTokens(tokenNames)

	contains := make(map[search.Token]bool, len(tokens))
	for _, token := range tokens {
		contains[token] = true
	}

	return func(q *search.Query) bool {
		e := q.Pat[len(q.Pat)-1]
		return !contains[e.Token]
	}
}
