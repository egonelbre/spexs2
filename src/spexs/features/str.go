package features

import (
	. "spexs"
)

var Str = [...]StrDesc{
	{"pat",
		"representation of the pattern",
		func(q *Query, db *Database) string {
			return q.String(db, true)
		}},
	{"pat-regexp",
		"representation of the pattern with group symbols replaced",
		func(q *Query, db *Database) string {
			return q.String(db, false)
		}},
}
