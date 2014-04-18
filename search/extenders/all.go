package extenders

import (
	"reflect"

	. "github.com/egonelbre/spexs2/search"
	"github.com/egonelbre/spexs2/utils"
)

type extender struct {
	Db *Database
}

type CreateFunc func(Setup, []byte) Extender

var All = [...]Extender{
	&Simple{},
	&Group{},
	&Star{},
	&StarGreedy{},
	&Regex{},
	&RegexGreedy{},
}

func mk(template Extender) Extender {
	ps := reflect.ValueOf(template)
	t := ps.Elem().Type()
	return reflect.New(t).Interface().(Extender)
}

func Make(template Extender, setup *Setup) (Extender, error) {
	v := mk(template)
	ps := reflect.ValueOf(v)
	s := ps.Elem()
	if s.Kind() == reflect.Struct {
		// assume "Db" is the first field
		field := s.FieldByName("Db")
		if field.IsValid() && field.CanSet() {
			field.Set(reflect.ValueOf(setup.Db))
		}
	}
	return v, nil
}

func Get(name string) (Extender, bool) {
	for _, fn := range All {
		if utils.UnqualifiedNameOf(fn) == name {
			return fn, true
		}
	}
	return nil, false
}

var Help = `
  Simple : uses the sequence tokens to discover the patterns  ( ACCT )
  Group : uses additionally defined groups in Alphabet.Groups ( AC[CT]T )
  Star : uses matching anything in the pattern ( AC.*T )
  StarGreedy : matches greedily anything in the pattern ( AC.*T )
  Regex : uses both group and star token in the pattern ( A[CT].*T )
  RegexGreedy : uses both group and star token in the pattern ( A[CT].*T )
`
