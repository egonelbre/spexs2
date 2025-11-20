package features

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/egonelbre/spexs2/search"
)

type CreateFunc any

type Desc struct {
	Category string
	Name     string
	Alias    []string
	Desc     string
	Create   CreateFunc
}

var All = [...]Desc{
	// only strings
	{
		Category: "Pattern",
		Name:     "Pat",
		Desc:     "Pat?()\tpattern as string",
		Create:   Pat,
	},
	{
		Category: "Pattern",
		Name:     "PatRegex",
		Alias:    []string{"Regex"},
		Desc:     "Regex?()\tpattern where groups have been expanded",
		Create:   PatRegex,
	},

	// pattern length related
	{
		Category: "Pattern",
		Name:     "PatLength",
		Desc:     "PatLength()\tpattern length",
		Create:   PatLength,
	},
	{
		Category: "Pattern",
		Name:     "PatChars",
		Desc:     "PatChars()\tcount of simple tokens in pattern",
		Create:   PatChars,
	},
	{
		Category: "Pattern",
		Name:     "PatGroups",
		Desc:     "PatGroups()\tcount of grouping tokens in pattern",
		Create:   PatGroups,
	},
	{
		Category: "Pattern",
		Name:     "PatStars",
		Desc:     "PatStars()\tcount of star tokens in pattern",
		Create:   PatStars,
	},

	// simple counting
	{
		Category: "Count",
		Name:     "Total",
		Desc:     "Total(A)\ttotal count of sequences",
		Create:   Total,
	},
	{
		Category: "Count",
		Name:     "Matches",
		Desc:     "Matches(A)\tcount of matching sequences",
		Create:   Matches,
	},
	{
		Category: "Count",
		Name:     "Seqs",
		Desc:     "Seqs(A)\tcount of unique sequences in matches",
		Create:   Seqs,
	},
	{
		Category: "Count",
		Name:     "Occs",
		Desc:     "Occs(A)\tcount of occurrences in the sequences",
		Create:   Occs,
	},

	// ratios and proportions
	{
		Category: "Ratio",
		Name:     "MatchesProp",
		Desc:     "MatchesProp(A)\t:= Matches(A)/Total(A)",
		Create:   MatchesProp,
	},
	{
		Category: "Ratio",
		Name:     "MatchesRatio",
		Desc:     "MatchesRatio(A, B)\t:= (Matches(A)+1)/(Matches(B)+1)",
		Create:   MatchesRatio,
	},
	{
		Category: "Ratio",
		Name:     "OccsRatio",
		Desc:     "OccsRatio(A, B)\t:= (Occs(A)+1)/(Occs(B)+1)",
		Create:   OccsRatio,
	},
	{
		Category: "Ratio",
		Name:     "MatchesPropRatio",
		Desc:     "MatchesPropRatio(A, B)\t:= ((Matches(A)+1)/(Total(A)+1))/((Matches(B)+1)/(Total(B)+1))",
		Create:   MatchesPropRatio,
	},

	// binomial
	{
		Category: "Statistics",
		Name:     "Binom",
		Desc:     "Binom(fore, back)\tbinomial p-value between fore and back datasets",
		Create:   Binom,
	},

	// hypergeometrics
	{
		Category: "Statistics",
		Name:     "Hyper",
		Desc:     "Hyper(fore, back)\thypergeometric p-value between for and back datsets",
		Create:   Hyper,
	},
	{
		Category: "Statistics",
		Name:     "HyperApprox",
		Desc:     "HyperApprox(fore, back)\tapprox. hypergeometric p-value (~5 sig. digits) between for and back datsets",
		Create:   HyperApprox,
	},
	{
		Category: "Statistics",
		Name:     "HyperDown",
		Desc:     "HyperDown(fore, back)\thypergeometric split down between for and back datsets",
		Create:   HyperDown,
	},
	{
		Category: "Statistics",
		Name:     "HyperOptimal",
		Desc:     "HyperOptimal(fore...)\tbest hypergeometric p-value across fore datasets",
		Create:   HyperOptimal,
	},
}

func Get(name string) (CreateFunc, bool) {
	for _, fn := range All {
		if fn.Name == name || slices.Contains(fn.Alias, name) {
			return fn.Create, true
		}
	}
	return nil, false
}

func CallCreateWithArgs(function CreateFunc, args []any) (search.Feature, error) {
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
	return inter.(search.Feature), nil
}

func Help() string {
	var b strings.Builder
	w := tabwriter.NewWriter(&b, 0, 0, 1, ' ', 0)
	fmt.Fprintf(w, "### Features\n")
	lastCategory := ""
	for _, fn := range All {
		if fn.Category != lastCategory {
			fmt.Fprintf(w, "\n  > %v\n", fn.Category)
			lastCategory = fn.Category
		}
		fmt.Fprintf(w, "  %v\n", fn.Desc)
	}
	_ = w.Flush()
	return b.String()
}

func functionAndType(fn any) (v reflect.Value, t reflect.Type, ok bool) {
	v = reflect.ValueOf(fn)
	ok = v.Kind() == reflect.Func
	if !ok {
		return
	}
	t = v.Type()
	return
}
