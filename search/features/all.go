package features

import (
	"reflect"

	. "github.com/egonelbre/spexs2/search"
	"github.com/egonelbre/spexs2/utils"
)

type feature struct {
	Db *Database
}

var All = [...]Feature{
	// simple counting
	&Total{}, &Matches{}, &Seqs{}, &Occs{},
	// ratios and proportions
	&MatchesProp{}, &MatchesRatio{}, &OccsRatio{}, &MatchesPropRatio{},
	// binomial
	&Binom{},
	// hypergeometrics
	&Hyper{}, &HyperApprox{}, &HyperOptimal{},
	// pattern length related
	&PatLength{}, &PatChars{}, &PatGroups{}, &PatStars{},
	// only strings
	&Pat{}, &PatRegex{},
}

func Get(name string) (Feature, bool) {
	for _, f := range All {
		if utils.UnqualifiedNameOf(f) == name {
			return f, true
		}
	}
	return nil, false
}

func mk(template Feature) Feature {
	ps := reflect.ValueOf(template)
	t := ps.Elem().Type()
	return reflect.New(t).Interface().(Feature)
}

func Make(template Feature, setup *Setup, args []interface{}) (Feature, error) {
	v := mk(template)
	ps := reflect.ValueOf(v)
	s := ps.Elem()
	if s.Kind() == reflect.Struct {
		// assume "Db" is the first field
		field := s.FieldByName("Db")
		index := 0
		if field.IsValid() && field.CanSet() {
			field.Set(reflect.ValueOf(setup.Db))
			index += 1
		}

		for _, v := range args {
			field := s.Field(index)
			index += 1
			if field.IsValid() && field.CanSet() {
				field.Set(reflect.ValueOf(v))
			}
		}
	}
	return v, nil
}

var Help = `
:Pattern:
  Pat?()   : pattern as string
  Regex?() : pattern where groups have been expanded

  PatLength() : pattern length
  PatChars()  : count of simple tokens in pattern
  PatGroups() : count of grouping tokens in pattern
  PatStars()  : count of star tokens in pattern

:Counting:
  // A and B refers to dataset
  Total(A)   : total count of sequences
  Matches(A) : count of matching sequences
  Seqs(A)    : count of unique sequences in matches
  Occs(A)    : count of occurences in the sequences

:Ratios:
  // A and B refer to datasets
  MatchesProp(A)     = Matches(A)/Total(A)
  MatchesRatio(A, B) = (Matches(A)+1)/(Matches(B)+1)
  OccsRatio(A, B)    = (Occs(A)+1)/(Occs(B)+1)
  MatchesPropRatio(A, B) = 
    ((Matches(A)+1)/(Total(A)+1))/((Matches(B)+1)/(Total(B)+1))

:Statistics:
  // fore and back refer to datasets
  Binom(fore, back)       : binomial p-value
  Hyper(fore, back)       : hypergeometric p-value
  HyperApprox(fore, back) : approx. hypergeometric p-value (~5 sig. digits)
  HyperDown(fore, back)   : hypergeometric split down
`

func functionAndType(fn interface{}) (v reflect.Value, t reflect.Type, ok bool) {
	v = reflect.ValueOf(fn)
	ok = v.Kind() == reflect.Func
	if !ok {
		return
	}
	t = v.Type()
	return
}
