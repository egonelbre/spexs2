package filters

import (
	. "spexs"
)

var All = [...]Desc{
	{"no-starting-group",
		"does not allow pattern to start with group",
		func(conf Conf) (Func, error) {
			return func(p *Query, ref *Database) bool {
				e := p.Pat[0]
				return e.IsGroup || e.IsStar
			}, nil
		}},
	{"no-ending-group",
		"does not allow pattern to end with group",
		func(conf Conf) (Func, error) {
			return func(p *Query, ref *Database) bool {
				e := p.Pat[len(p.Pat)-1]
				return !e.IsGroup
			}, nil
		}},
}
