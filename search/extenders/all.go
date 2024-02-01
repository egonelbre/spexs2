package extenders

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/egonelbre/spexs2/search"
)

type CreateFunc func(search.Setup, []byte) search.Extender

type Desc struct {
	Name string
	Desc string
	Func search.Extender
}

var All = [...]Desc{
	{
		Name: "Simple",
		Desc: "uses the sequence tokens to discover the patterns\t( ACCT )",
		Func: Simple,
	},
	{
		Name: "Group",
		Desc: "uses additionally defined groups in Alphabet.Groups\t( AC[CT]T )",
		Func: Group,
	},
	{
		Name: "Star",
		Desc: "uses matching anything in the pattern\t( AC.*T )",
		Func: Star,
	},
	{
		Name: "StarGreedy",
		Desc: "matches greedily anything in the pattern\t( AC.*T )",
		Func: StarGreedy,
	},
	{
		Name: "Regex",
		Desc: "uses both group and star token in the pattern\t( A[CT].*T )",
		Func: Regex,
	},
	{
		Name: "RegexGreedy",
		Desc: "uses both group and star token in the pattern\t( A[CT].*T )",
		Func: RegexGreedy,
	},
}

func wrap(f search.Extender) CreateFunc {
	return func(s search.Setup, data []byte) search.Extender {
		return f
	}
}

func Get(name string) (search.Extender, bool) {
	for _, fn := range All {
		if fn.Name == name {
			return fn.Func, true
		}
	}
	return nil, false
}

func Help() string {
	var b strings.Builder
	w := tabwriter.NewWriter(&b, 0, 0, 1, ' ', 0)
	fmt.Fprintf(w, "### Extenders\n")
	for _, fn := range All {
		fmt.Fprintf(w, "  %v\t%v\n", fn.Name, fn.Desc)
	}
	w.Flush()
	return b.String()
}
