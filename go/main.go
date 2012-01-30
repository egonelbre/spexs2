package main

import "./spexs"

func main() {
	input := NewFifoPool()
	results := NewFifoPool()
	ref := &UnicodeReference{}
	RunParallel(ref, input, results, GroupStarExtender, TrieCountFilter(7))
	for p, valid := results.Take(); valid {
		print p.ToString()
	}
}