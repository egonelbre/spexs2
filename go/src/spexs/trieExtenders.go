package spexs

func output(out Patterns, patterns map[Char]*TrieNode) {
	for _, node := range patterns {
		out <- node
	}
}

func trieSimpleExtend(n *TrieNode, ref Reference, patterns map[Char]*TrieNode) {
	indices, poss := n.Pos.Iter()
	for idx := range indices {
		mpos := <- poss
		plen := byte(len(ref.(*UnicodeReference).Pats[idx].Pat) - 1)
		var k byte
		for k = 0; (k < plen) && (mpos > 0); k += 1 {
			if mpos & (1 << k) == 0 { continue }
			mpos &^= 1 << k

			char, next, _ := ref.Next(idx, k)
			pat, exists := patterns[char]
			if !exists {
				pat = NewTrieNode(char, n)
				patterns[char] = pat
			}
			pat.Pos.Add(idx, next)
		}
	}
}

func SimpleExtender(p Pattern, ref Reference) Patterns {
	result := MakePatterns()
	patterns := make(map[Char]*TrieNode)
	
	node := p.(*TrieNode)
	trieSimpleExtend(node, ref, patterns)

	output(result, patterns)
	close(result)
	return result
}


func trieGroupCombine(n *TrieNode, ref Reference, patterns map[Char]*TrieNode) {
	for _, g := range ref.(*UnicodeReference).Groups {
		pat := NewTrieNode(g.Id, n)
		pat.IsGroup = true
		patterns[g.Id] = pat
		for _, char := range g.Chars {
			if _, exists := patterns[char]; exists {
				SetAddSet(patterns[g.Id].Pos, patterns[char].Pos)
			}
		}
	}
}

func GroupExtender(p Pattern, ref Reference) Patterns {
	result := MakePatterns()
	patterns := make(map[Char]*TrieNode)

	node := p.(*TrieNode)
	trieSimpleExtend(node, ref, patterns)
	trieGroupCombine(node, ref, patterns)

	output(result, patterns)
	close(result)
	return result
}

func trieStarExtend(node *TrieNode, ref Reference, stars map[Char]*TrieNode) {
	indices, poss := node.Pos.Iter()
	for idx := range indices {
		mpos := <- poss
		plen := byte(len(ref.(*UnicodeReference).Pats[idx].Pat) - 1)
		if mpos == 0 { continue }

		var k byte
		for k = 0; k < plen; k += 1 {
			if mpos & (1 << k) == 0 { break }
		}

		char, next, valid := ref.Next(idx, k)
		for valid {
			pat, exists := stars[char]
			if !exists {
				pat = NewTrieNode(char, node)
				pat.IsStar = true
				stars[char] = pat
			}
			pat.Pos.Add(idx, next)
			char, next, valid = ref.Next(idx, next)
		}
	}
}

func StarExtender(p Pattern, ref Reference) Patterns {
	result := MakePatterns()
	patterns := make(map[Char]*TrieNode)
	stars := make(map[Char]*TrieNode)

	node := p.(*TrieNode)
	trieSimpleExtend(node, ref, patterns)
	trieStarExtend(node, ref, stars)

	output(result, patterns)
	output(result, stars)
	close(result)
	return result
}

func GroupStarExtender(p Pattern, ref Reference) Patterns {
	result := MakePatterns()
	patterns := make(map[Char]*TrieNode)
	stars := make(map[Char]*TrieNode)

	node := p.(*TrieNode)
	trieSimpleExtend(node, ref, patterns)
	trieGroupCombine(node, ref, patterns)
	trieStarExtend(node, ref, stars)
	trieGroupCombine(node, ref, stars)

	output(result, patterns)
	output(result, stars)

	close(result)
	return result
}
