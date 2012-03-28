package spexs

func output(out TrieNodes, patterns map[Char]*TrieNode) {
	for _, node := range patterns {
		out <- node
	}
}

func trieSimpleExtend(node *TrieNode, ref *UnicodeReference, 
		patterns map[Char]*TrieNode) {
	indices, poss := node.Pos.Iter()
	for idx := range indices {
		mpos := <- poss
		plen := byte(len(ref.Pats[idx].Pat))
		var k byte
		for k = 0; (k < plen) && (mpos > 0); k += 1 {
			if mpos & (1 << k) == 0 { continue }
			mpos &^= 1 << k

			char, next, valid := ref.Next(idx, k)
			if !valid { break }

			pat, exists := patterns[char]
			if !exists {
				pat = NewTrieNode(char, node)
				patterns[char] = pat
			}
			pat.Pos.Add(idx, next)
		}
	}
}

func SimpleExtender(node *TrieNode, ref *UnicodeReference) TrieNodes {
	result := MakeTrieNodes()
	patterns := make(map[Char]*TrieNode)
	
	trieSimpleExtend(node, ref, patterns)

	output(result, patterns)
	close(result)
	return result
}


func trieGroupCombine(node *TrieNode, ref *UnicodeReference, patterns map[Char]*TrieNode, star bool) {
	for _, g := range ref.Groups {
		pat := NewTrieNode(g.Id, node)
		pat.IsGroup = true
		pat.IsStar = star
		patterns[g.Id] = pat
		for _, char := range g.Chars {
			if _, exists := patterns[char]; exists {
				SetAddSet(patterns[g.Id].Pos, patterns[char].Pos)
			}
		}
	}
}

func GroupExtender(node *TrieNode, ref *UnicodeReference) TrieNodes {
	result := MakeTrieNodes()
	patterns := make(map[Char]*TrieNode)

	trieSimpleExtend(node, ref, patterns)
	trieGroupCombine(node, ref, patterns, false)

	output(result, patterns)
	close(result)
	return result
}

func trieStarExtend(node *TrieNode, ref Reference, stars map[Char]*TrieNode) {
	indices, poss := node.Pos.Iter()
	for idx := range indices {
		mpos := <- poss
		plen := byte(len(ref.(*UnicodeReference).Pats[idx].Pat))
		if mpos == 0 { continue }

		var k byte
		for k = 0; k < plen; k += 1 {
			if mpos & (1 << k) != 0 { break }
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

func StarExtender(node *TrieNode, ref *UnicodeReference) TrieNodes {
	result := MakeTrieNodes()
	patterns := make(map[Char]*TrieNode)
	stars := make(map[Char]*TrieNode)
	trieSimpleExtend(node, ref, patterns)
	trieStarExtend(node, ref, stars)

	output(result, patterns)
	output(result, stars)
	close(result)
	return result
}

func GroupStarExtender(node *TrieNode, ref *UnicodeReference) TrieNodes {
	result := MakeTrieNodes()
	patterns := make(map[Char]*TrieNode)
	stars := make(map[Char]*TrieNode)

	trieSimpleExtend(node, ref, patterns)
	trieGroupCombine(node, ref, patterns, false)
	trieStarExtend(node, ref, stars)
	trieGroupCombine(node, ref, stars, true)

	output(result, patterns)
	output(result, stars)

	close(result)
	return result
}
