package features

import (
	. "spexs"
	"stats/hyper"
)

var All = [...]Desc{
	{"query-seqs",
		"total number of query sequences",
		func(p *Pattern, ref *Reference) float64 {
			return float64(ref.Groupings[0])
		}},
	{"back-seqs",
		"total number of background sequences",
		func(p *Pattern, ref *Reference) float64 {
			return float64(ref.Groupings[1])
		}},
	{"query-match-seqs",
		"number of matching query sequences",
		func(p *Pattern, ref *Reference) float64 {
			return float64(p.Count(ref, 0))
		}},
	{"back-match-seqs",
		"number of matching background sequences",
		func(p *Pattern, ref *Reference) float64 {
			return float64(p.Count(ref, 1))
		}},
	{"query-match-occs",
		"number of occurences in query",
		func(p *Pattern, ref *Reference) float64 {
			return float64(p.Occs(ref, 0))
		}},
	{"back-match-occs",
		"number of occurences in background",
		func(p *Pattern, ref *Reference) float64 {
			return float64(p.Occs(ref, 1))
		}},
	{"query-match-seqs-prop",
		"percentage of matching sequences in query",
		func(p *Pattern, ref *Reference) float64 {
			return float64(p.Count(ref, 0)) / float64(ref.Groupings[0])
		}},
	{"back-match-seqs-prop",
		"percentage of matching sequences in background",
		func(p *Pattern, ref *Reference) float64 {
			return float64(p.Count(ref, 1)) / float64(ref.Groupings[1])
		}},

	{"match-hyper-up-pvalue",
		"hypergeometric split p-value",
		func(p *Pattern, ref *Reference) float64 {
			query := p.Count(ref, 0)
			back := p.Count(ref, 1)
			pvalue := hyper.Split(query, back, ref.Groupings[0], ref.Groupings[1])
			return pvalue
		}},
	{"match-hyper-up-pvalue-approx",
		"approximate hypergeometric split p-value (~5 significant digits)",
		func(p *Pattern, ref *Reference) float64 {
			query := p.Count(ref, 0)
			back := p.Count(ref, 1)
			pvalue := hyper.SplitApprox(query, back, ref.Groupings[0], ref.Groupings[1])
			return pvalue
		}},
	{"match-hyper-down-pvalue",
		"hypergeometric split p-value down",
		func(p *Pattern, ref *Reference) float64 {
			query := p.Count(ref, 0)
			back := p.Count(ref, 1)
			pvalue := hyper.SplitDown(query, back, ref.Groupings[0], ref.Groupings[1])
			return pvalue
		}},
	{"match-ratio",
		"ratio of (matches in query + 1) / (matches in background + 1)",
		func(p *Pattern, ref *Reference) float64 {
			query := p.Count(ref, 0)
			back := p.Count(ref, 1)
			return float64(query+1) / float64(back+1)
		}},

	{"pat-length",
		"length of the pattern",
		func(p *Pattern, ref *Reference) float64 {
			t := 0
			for _, e := range p.Pat {
				t += 1
				if e.IsStar {
					t += 1
				}
			}
			return float64(t)
		}},
	{"pat-chars",
		"count of characters in pattern",
		func(p *Pattern, ref *Reference) float64 {
			t := 0
			for _, e := range p.Pat {
				if !e.IsGroup {
					t += 1
				}
			}
			return float64(t)
		}},
	{"pat-groups",
		"count of groups in pattern",
		func(p *Pattern, ref *Reference) float64 {
			t := 0
			for _, e := range p.Pat {
				if e.IsGroup {
					t += 1
				}
			}
			return float64(t)
		}},
	{"pat-stars",
		"count of stars in pattern",
		func(p *Pattern, ref *Reference) float64 {
			t := 0
			for _, e := range p.Pat {
				if e.IsStar {
					t += 1
				}
			}
			return float64(t)
		}},
}
