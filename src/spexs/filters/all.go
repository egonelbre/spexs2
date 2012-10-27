package filters

import (
	. "spexs"
	"strings"
	"utils"
)

var All = [...]Desc{
	{"no-starting-group",
		"does not allow pattern to start with group",
		func(conf Conf, setup Setup) (Func, error) {
			return func(q *Query) bool {
				e := q.Pat[0]
				return !(e.IsGroup || e.IsStar)
			}, nil
		}},
	{"no-ending-group",
		"does not allow pattern to end with group",
		func(conf Conf, setup Setup) (Func, error) {
			return func(q *Query) bool {
				e := q.Pat[len(q.Pat)-1]
				return !e.IsGroup
			}, nil
		}},
	{"no-ending-tokens",
		"does not allow pattern to end with token",
		func(conf Conf, setup Setup) (Func, error) {
			var filt struct{ Tokens string }
			utils.ApplyObject(&conf, &filt)

			line := strings.TrimSpace(filt.Tokens)
			tokenNames := strings.Split(line, setup.Db.Separator)
			tokens := setup.Db.ToTokens(tokenNames)

			contains := make(map[Token]bool, len(tokens))
			for _, token := range tokens {
				contains[token] = true
			}

			return func(q *Query) bool {
				e := q.Pat[len(q.Pat)-1]
				return !contains[e.Token]
			}, nil
		}},
}
