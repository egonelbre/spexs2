package trie

import "spexs"

func output(out Patterns, patterns map[Char]Pattern) {
	for _, node := range patterns {
		out <- node
	}
}

func simpleExtend(node Pattern, ref Reference, patterns map[Char]Pattern) {

	for idx, mpos := range node.Pos.Iter() {
		plen := byte(len(ref.Pats[idx].Pat))
		var k byte
		for k = 0; (k < plen) && (mpos > 0); k += 1 {
			if mpos&(1<<k) == 0 {
				continue
			}
			mpos &^= 1 << k

			char, next, valid := ref.Next(idx, k)
			if !valid {
				break
			}

			pat, exists := patterns[char]
			if !exists {
				pat = NewNode(char, node)
				patterns[char] = pat
			}
			pat.Pos.Add(idx, next)
		}
	}
}

func SimpleExtender(node Pattern, ref Reference) Patterns {
	result := NewPatterns()
	patterns := make(map[Char]Pattern)

	simpleExtend(node, ref, patterns)

	output(result, patterns)
	close(result)
	return result
}

func groupCombine(node Pattern, ref Reference, patterns map[Char]Pattern, star bool) {
	for _, g := range ref.Groups {
		pat := NewNode(g.Id, node)
		pat.IsGroup = true
		pat.IsStar = star
		patterns[g.Id] = pat
		for _, char := range g.Chars {
			if _, exists := patterns[char]; exists {
				spexs.SetAddSet(patterns[g.Id].Pos, patterns[char].Pos)
			}
		}
	}
}

func GroupExtender(node Pattern, ref Reference) Patterns {
	result := NewPatterns()
	patterns := make(map[Char]Pattern)

	simpleExtend(node, ref, patterns)
	groupCombine(node, ref, patterns, false)

	output(result, patterns)
	close(result)
	return result
}

func trieStarExtend(node Pattern, ref Reference, stars map[Char]Pattern) {
	for idx, mpos := range node.Pos.Iter() {
		plen := byte(len(ref.Pats[idx].Pat))
		if mpos == 0 {
			continue
		}

		var k byte
		for k = 0; k < plen; k += 1 {
			if mpos&(1<<k) != 0 {
				break
			}
		}

		char, next, valid := ref.Next(idx, k)
		for valid {
			pat, exists := stars[char]
			if !exists {
				pat = NewNode(char, node)
				pat.IsStar = true
				stars[char] = pat
			}
			pat.Pos.Add(idx, next)
			char, next, valid = ref.Next(idx, next)
		}
	}
}

func StarExtender(node Pattern, ref Reference) Patterns {
	result := NewPatterns()
	patterns := make(map[Char]Pattern)
	stars := make(map[Char]Pattern)
	simpleExtend(node, ref, patterns)
	trieStarExtend(node, ref, stars)

	output(result, patterns)
	output(result, stars)
	
	close(result)
	return result
}

func GroupStarExtender(node Pattern, ref Reference) Patterns {
	result := NewPatterns()
	patterns := make(map[Char]Pattern)
	stars := make(map[Char]Pattern)

	simpleExtend(node, ref, patterns)
	groupCombine(node, ref, patterns, false)
	trieStarExtend(node, ref, stars)
	groupCombine(node, ref, stars, true)

	output(result, patterns)
	output(result, stars)

	close(result)
	return result
}
