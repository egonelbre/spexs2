package features

import (
	. "spexs"
)

var Str = [...]StrDesc{
	{"pat",
		"representation of the pattern",
		func(q *Query) string {
			return q.String(true)
		}},
	{"pat-regexp",
		"representation of the pattern with group symbols replaced",
		func(q *Query) string {
			return q.String(false)
		}},
}
