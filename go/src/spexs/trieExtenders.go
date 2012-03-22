package spexs

func output(out Patterns, patterns map[Char]*TrieNode) {
	for _, node := range patterns {
		out <- *node
	}
}

func trieSimpleExtend(n *TrieNode, ref Reference, patterns map[Char]*TrieNode) {
	for pos := range n.Pos.Iter() {
		char, next, valid := ref.Next(pos)
		if !valid {
			continue
		}
		pat, exists := patterns[char]
		if !exists {
			pat = NewTrieNode(char, n)
			patterns[char] = pat
		}
		pat.Pos.Add(next)
	}
}

func SimpleExtender(p Pattern, ref Reference) Patterns {
	result := MakePatterns()
	patterns := make(map[Char]*TrieNode)
	
	node := p.(TrieNode)
	trieSimpleExtend(&node, ref, patterns)

	output(result, patterns)
	close(result)
	return result
}

func trieGroupCombine(n *TrieNode, ref Reference, patterns map[Char]*TrieNode) {
	for _, g := range ref.(UnicodeReference).Groups {
		pat := NewTrieNode(g.Id, n)
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

	node := p.(TrieNode)
	trieSimpleExtend(&node, ref, patterns)
	trieGroupCombine(&node, ref, patterns)

	output(result, patterns)
	close(result)
	return result
}

func trieStarExtend(node *TrieNode, ref Reference, stars map[Char]*TrieNode) {
	lowest := make(map[int]Pos)

	for p := range node.Pos.Iter() {
		idx, _ := PosDecode(p)
		cur, exists := lowest[idx]
		if !exists || p < cur {
			lowest[idx] = p
		}
	}

	for _, p := range lowest {
		char, next, valid := ref.Next(p);
		for valid {
			pat, exists := stars[char]
			if !exists {
				pat = NewTrieNode(char, node)
				pat.IsStar = true
				stars[char] = pat
			}
			pat.Pos.Add(next)
			char, next, valid = ref.Next(next)
		}
	}
}

func StarExtender(p Pattern, ref Reference) Patterns {
	result := MakePatterns()
	patterns := make(map[Char]*TrieNode)
	stars := make(map[Char]*TrieNode)

	node := p.(TrieNode)
	trieSimpleExtend(&node, ref, patterns)
	trieStarExtend(&node, ref, stars)

	output(result, patterns)
	output(result, stars)
	close(result)
	return result
}

func GroupStarExtender(p Pattern, ref Reference) Patterns {
	result := MakePatterns()
	patterns := make(map[Char]*TrieNode)
	stars := make(map[Char]*TrieNode)

	node := p.(TrieNode)
	trieSimpleExtend(&node, ref, patterns)
	trieGroupCombine(&node, ref, patterns)
	trieStarExtend(&node, ref, stars)
	trieGroupCombine(&node, ref, stars)

	output(result, patterns)
	output(result, stars)

	close(result)
	return result
}
