package main

import (
	"stats"
	. "spexs/trie"
)

type FeatureFunc func(*Pattern, *Reference) float64

var Features = map[string]FeatureFunc{
	"length": func(p *Pattern, ref *Reference) float64 {
		return float64(p.Len())
	},
	"ng": func(p *Pattern, ref *Reference) float64 {
		t := 0
		for p != nil {
			if ! (p.IsGroup || p.IsStar) {
				t += 1
			}
			p = p.Parent
		}
		return float64(t)
	},
	"count": func(p *Pattern, ref *Reference) float64 {
		return float64(p.Pos.Len()) / float64(ref.Groupings[0])
	},
	"p-value": func(p *Pattern, ref *Reference) float64 {
		c := p.Count(ref)
		pvalue := stats.HypergeometricSplit(c[0], c[1], ref.Groupings[0], ref.Groupings[1])
		return pvalue
	},
}
