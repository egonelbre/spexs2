package main

import . "spexs"

type TrieFilterCreator func(interface{}) TrieFilterFunc

var filters = map[string] PatternFilterCreator {
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

func CreateFilter(conf map[string]interface{}, setup Setup) TrieFilterFunc {
	
}