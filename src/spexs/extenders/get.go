package extenders

import (
	. "spexs"
	"utils"
)

func Get(name string) (Extender, bool) {
	for _, fn := range All {
		if utils.FuncName(fn) == name {
			return fn, true
		}
	}
	return nil, false
}
