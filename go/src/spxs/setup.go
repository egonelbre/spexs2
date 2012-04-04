package main

type Setup struct {
	Ref        *UnicodeReference
	Out        TriePooler
	In         TriePooler
	Extender   TrieExtenderFunc
	Extendable TrieFilterFunc
	Outputtable TrieFilterFunc
}

var extenders = map[string]TrieExtenderFunc{
	"simple": SimpleExtender,
	"group":  GroupExtender,
	"star":   StarExtender,
	"regexp": GroupStarExtender,
}

type PatternFilterCreator func(limit int) TrieFilterFunc

var limiters = map[string]PatternFilterCreator{
	"count": func(limit int) TrieFilterFunc {
		return func(p *TrieNode) bool {
			return p.Pos.Len() >= limit
		}
	},
	"length": func(limit int) TrieFilterFunc {
		return func(p *TrieNode) bool {
			return p.Len() <= limit
		}
	},
	"complexity": func(limit int) TrieFilterFunc {
		return func(p *TrieNode) bool {
			return p.Complexity() <= limit
		}
	},
}

var fitnesses = map[string]TrieFitnessFunc{
	"def": func(p *TrieNode) float64 {
		return float64(p.Len() * p.Pos.Len())
	},
	"len": func(p *TrieNode) float64 {
		return float64(p.Len())
	},
	"count": func(p *TrieNode) float64 {
		return float64(p.Pos.Len())
	},
	"complexity": func(p *TrieNode) float64 {
		return float64(p.Complexity())
	},
}

func inputOrdering(p *TrieNode) float64 {
	return 1 / float64(p.Len())
}
