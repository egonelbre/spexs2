package filters

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/egonelbre/spexs2/search"
)

type CreateFunc func(*search.Setup, []byte) search.Filter

type Desc struct {
	Category string
	Name     string
	Desc     string
	Create   CreateFunc
}

var All = [...]Desc{
	{
		Category: "Pattern",
		Name:     "NoStartingGroup",
		Desc:     "removes patterns with starting group token",
		Create:   NoStartingGroup,
	},
	{
		Category: "Pattern",
		Name:     "NoEndingGroup",
		Desc:     "removes patterns with ending group token",
		Create:   NoEndingGroup,
	},
	{
		Category: "Pattern",
		Name:     "NoTokens",
		Desc:     "removes patterns ending with tokens specified in \"Tokens\" argument\n\t(useful only in output.filter)",
		Create:   NoTokens,
	},
}

func Get(name string) (CreateFunc, bool) {
	for _, fn := range All {
		if fn.Name == name {
			return fn.Create, true
		}
	}
	return nil, false
}

func Help() string {
	var b strings.Builder
	w := tabwriter.NewWriter(&b, 0, 0, 1, ' ', 0)
	fmt.Fprintf(w, "### Filters\n")
	lastCategory := ""
	for _, fn := range All {
		if fn.Category != lastCategory {
			fmt.Fprintf(w, "\n  > %v\n", fn.Category)
			lastCategory = fn.Category
		}
		fmt.Fprintf(w, "  %v\t%v\n", fn.Name, fn.Desc)
	}
	fmt.Fprintf(w, "\n  > Feature\n")
	fmt.Fprintf(w, "  Any feature can be used as a filter.\n")
	_ = w.Flush()
	return b.String()
}
