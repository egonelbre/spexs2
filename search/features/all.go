package features

import (
	"fmt"
	"reflect"

	. "github.com/egonelbre/spexs2/search"
	"github.com/egonelbre/spexs2/utils"
)

type CreateFunc interface{}

var All = [...]CreateFunc{
	// simple counting
	Total, Matches, Seqs, Occs,
	// ratios and proportions
	MatchesProp, MatchesRatio, OccsRatio, MatchesPropRatio,
	// binomial
	Binom,
	// hypergeometrics
	Hyper, HyperApprox, HyperDown, HyperOptimal,
	// pattern length related
	PatLength, PatChars, PatGroups, PatStars,
	// only strings
	Pat, PatRegex,
}

func Get(name string) (CreateFunc, bool) {
	for _, fn := range All {
		if utils.FuncName(fn) == name {
			return fn, true
		}
	}
	return nil, false
}

func CallCreateWithArgs(function CreateFunc, args []interface{}) (Feature, error) {
	fn, fnType, ok := functionAndType(function)
	if !ok {
		return nil, fmt.Errorf("argument is not a function")
	}

	if fnType.NumIn() != len(args) {
		return nil, fmt.Errorf("invalid number of arguments, requires %v", fnType.NumIn())
	}

	arguments := make([]reflect.Value, fnType.NumIn())
	for i := range args {
		arguments[i] = reflect.ValueOf(args[i])
	}
	result := fn.Call(arguments)
	inter := result[0].Interface()
	return inter.(Feature), nil
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
