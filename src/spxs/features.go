package main

import (
	. "spexs/trie"
	"stats"
)

type Feature struct {
	Desc string
	Func FeatureFunc
}

type FeatureFunc func(*Pattern, *Reference) float64

var Features = map[string]FeatureFunc{
	"query-seqs": func(p *Pattern, ref *Reference) float64 {
		return float64(ref.Groupings[0])
	},
	"back-seqs": func(p *Pattern, ref *Reference) float64 {
		return float64(ref.Groupings[1])
	},

	"query-match-seqs": func(p *Pattern, ref *Reference) float64 {
		return float64(p.Count(ref,0))
	},
	"back-match-seqs": func(p *Pattern, ref *Reference) float64 {
		return float64(p.Count(ref,1))
	},
	"query-match-occs": func(p *Pattern, ref *Reference) float64 {
		return float64(p.Occs(ref,0))
	},
	"back-match-occs": func(p *Pattern, ref *Reference) float64 {
		return float64(p.Occs(ref,1))
	},
	"query-match-seqs-prop": func(p *Pattern, ref *Reference) float64 {
		return float64(p.Count(ref,0)/ref.Groupings[0])
	},
	"back-match-seqs-prop": func(p *Pattern, ref *Reference) float64 {
		return float64(p.Count(ref,1)/ref.Groupings[1])
	},
	
	"match-hyper-up-pvalue": func(p *Pattern, ref *Reference) float64 {
		query := p.Count(ref,0)
		back := p.Count(ref,1)
		pvalue := stats.HypergeometricSplit(query, back, ref.Groupings[0], ref.Groupings[1])
		return pvalue
	},
	"match-hyper-down-pvalue": func(p *Pattern, ref *Reference) float64 {
		query := p.Count(ref,0)
		back := p.Count(ref,1)
		pvalue := stats.HypergeometricSplitDown(query, back, ref.Groupings[0], ref.Groupings[1])
		return pvalue
	},
	"match-ratio": func(p *Pattern, ref *Reference) float64 {
		query := p.Count(ref,0)
		back := p.Count(ref,1)
		return float64((query+1)/(back+1))
	},

	"pat-length": func(p *Pattern, ref *Reference) float64 {
		t := 0
		for p != nil {
			t += 1
			p = p.Parent
		}
		return float64(t)
	},
	"pat-chars": func(p *Pattern, ref *Reference) float64 {
		t := 0
		for p != nil {
			if !p.IsGroup {
				t += 1
			}
			p = p.Parent
		}
		return float64(t)
	},
	"pat-groups": func(p *Pattern, ref *Reference) float64 {
		t := 0
		for p != nil {
			if p.IsGroup {
				t += 1
			}
			p = p.Parent
		}
		return float64(t)
	},
	"pat-stars": func(p *Pattern, ref *Reference) float64 {
		t := 0
		for p != nil {
			if p.IsStar {
				t += 1
			}
			p = p.Parent
		}
		return float64(t)
	},
}
