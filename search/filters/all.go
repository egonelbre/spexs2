package filters

import (
	"reflect"

	. "github.com/egonelbre/spexs2/search"
	"github.com/egonelbre/spexs2/utils"
)

type filter struct {
	Db *Database
}

var All = [...]Filter{
	&NoStartingGroup{},
	&NoEndingGroup{},
	&NoTokens{},
}

func Get(name string) (Filter, bool) {
	for _, f := range All {
		if utils.UnqualifiedNameOf(f) == name {
			return f, true
		}
	}
	return nil, false
}

func mk(template Filter) Filter {
	ps := reflect.ValueOf(template)
	t := ps.Elem().Type()
	return reflect.New(t).Interface().(Filter)
}

func trySetField(obj interface{}, field string, value interface{}) {
	defer recover()
	// pointer to struct
	ps := reflect.ValueOf(obj)
	// struct itself
	s := ps.Elem()
	if s.Kind() == reflect.Struct {
		field := s.FieldByName(field)
		if field.IsValid() && field.CanSet() {
			field.Set(reflect.ValueOf(value))
		}
	}
}

func Make(template Filter, setup *Setup, data []byte) Filter {
	v := mk(template)
	trySetField(v, "Db", setup.Db)
	return v
}

var Help = `
:Pattern:
  NoStartingGroup() : removes patterns with starting group token
  NoEndingGroup()   : removes patterns with ending group token
  	                (useful only in output.filter)

  NoTokens() : removes patterns ending with tokens specified in "Tokens" argument

:Feature:
  Any feature can be used as a filter
`
