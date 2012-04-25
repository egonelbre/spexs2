package main

import (
	. "spexs/trie"
)

type FeatureFunc func(*Pattern, *Reference) float64

var Features = map[string]FeatureFunc{
	"length": func(p *Pattern, ref *Reference) float64 {
		return float64(p.Len())
	},
	"ng": func(p *Pattern, ref *Reference) float64 {
		return float64(p.NG())
	},
	"count": func(p *Pattern, ref *Reference) float64 {
		return float64(p.Pos.Len()) / float64(ref.Groupings[0])
	},
	"p-value": func(p *Pattern, ref *Reference) float64 {
		return p.PValue(ref)
	},
}
