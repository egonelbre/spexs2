package features

import (
	. "spexs"
)

type Func func(*Query, *Database) float64

type Desc struct {
	Name string
	Desc string
	Func Func
}

func Get(name string) (*Desc, bool) {
	for _, e := range All {
		if e.Name == name {
			return &e, true
		}
	}
	return nil, false
}

type StrFunc func(*Query, *Database) string

type StrDesc struct {
	Name string
	Desc string
	Func StrFunc
}

func GetStr(name string) (*StrDesc, bool) {
	for _, e := range Str {
		if e.Name == name {
			return &e, true
		}
	}
	return nil, false
}
