package features

import (
	. "spexs"
	"utils"
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
			pvalue := utils.HypergeometricSplit(query, back, ref.Groupings[0], ref.Groupings[1])
			return pvalue
		}},
	{"match-hyper-up-pvalue-approx",
		"hypergeometric split p-value (approximate)",
		func(p *Pattern, ref *Reference) float64 {
			query := p.Count(ref, 0)
			back := p.Count(ref, 1)
			pvalue := utils.HypergeometricSplitApprox(query, back, ref.Groupings[0], ref.Groupings[1])
			return pvalue
		}},
	{"match-hyper-down-pvalue",
		"hypergeometric split p-value down",
		func(p *Pattern, ref *Reference) float64 {
			query := p.Count(ref, 0)
			back := p.Count(ref, 1)
			pvalue := utils.HypergeometricSplitDown(query, back, ref.Groupings[0], ref.Groupings[1])
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
			t := -1 // because first "" is also a char
			for p != nil {
				t += 1
				if p.IsStar {
					t += 1
				}
				p = p.Parent
			}
			return float64(t)
		}},
	{"pat-chars",
		"count of characters in pattern",
		func(p *Pattern, ref *Reference) float64 {
			t := -1 // because first "" is also a char
			for p != nil {
				if !p.IsGroup {
					t += 1
				}
				p = p.Parent
			}
			return float64(t)
		}},
	{"pat-groups",
		"count of groups in pattern",
		func(p *Pattern, ref *Reference) float64 {
			t := 0
			for p != nil {
				if p.IsGroup {
					t += 1
				}
				p = p.Parent
			}
			return float64(t)
		}},
	{"pat-stars",
		"count of stars in pattern",
		func(p *Pattern, ref *Reference) float64 {
			t := 0
			for p != nil {
				if p.IsStar {
					t += 1
				}
				p = p.Parent
			}
			return float64(t)
		}},
}
