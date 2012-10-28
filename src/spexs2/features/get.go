package features

import "utils"

func Get(name string) (interface{}, bool) {
	for _, fn := range All {
		if utils.FuncName(fn) == name {
			return fn, true
		}
	}
	return nil, false
}
