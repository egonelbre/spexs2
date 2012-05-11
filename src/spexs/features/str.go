package features

import (
	. "spexs"
)

var Str = [...]StrDesc{
	{"pat",
		"representation of the pattern",
		func(p *Pattern, ref *Reference) string {
			return p.String()
		}},
	{"pat-regexp",
		"representation of the pattern with group symbols replaced",
		func(p *Pattern, ref *Reference) string {
			return ref.ReplaceGroups(p.String())
		}},
}
