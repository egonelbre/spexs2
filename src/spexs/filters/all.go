package filters

import (
	. "spexs"
)

var All = [...]Desc{
	{"no-starting-group",
		"does not allow pattern to start with group",
		func(conf Conf) (Func, error) {
			return func(p *Pattern, ref *Reference) bool {
				for p != nil {
					if (p.IsGroup || p.IsStar) &&
						(p.Parent != nil) &&
						(p.Parent.Parent == nil) {
						return false
					}
					p = p.Parent
				}
				return true
			}, nil
		}},
	{"no-ending-group",
		"does not allow pattern to end with group",
		func(conf Conf) (Func, error) {
			return func(p *Pattern, ref *Reference) bool {
				return !p.IsGroup
			}, nil
		}},
}
