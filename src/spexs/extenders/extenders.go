package extenders

import (
	. "spexs"
	"utils"
)

func output(out Patterns, patterns map[Char]*Pattern) {
	for _, node := range patterns {
		out <- node
	}
}

func extend(node *Pattern, ref *Reference, patterns map[Char]*Pattern) {
	for idx, ipos := range node.Pos.Iter() {
		px := ref.Seqs[idx].Pat
		plen := uint(len(px))
		for k := uint(0); k < plen; k += 1 {
			if (ipos >> k & 1) == 0 {
				continue
			}
			char := Char(px[k])
			pat, exists := patterns[char]
			if !exists {
				pat = NewPattern(char, node, false, false)
				patterns[char] = pat
			}
			pat.Pos.Add(idx, k+1)
		}
	}
}

func Simplex(node *Pattern, ref *Reference) Patterns {
	result := NewPatterns()
	patterns := make(map[Char]*Pattern)

	extend(node, ref, patterns)

	output(result, patterns)
	close(result)
	return result
}

func combine(node *Pattern, ref *Reference, patterns map[Char]*Pattern, star bool) {
	for _, g := range ref.Groups {
		pat := NewPattern(g.Id, node, true, star)
		patterns[g.Id] = pat
		for _, char := range g.Chars {
			if _, exists := patterns[char]; exists {
				patterns[g.Id].Pos.AddSet(patterns[char].Pos)
			}
		}
	}
}

func Groupex(node *Pattern, ref *Reference) Patterns {
	result := NewPatterns()
	patterns := make(map[Char]*Pattern)

	extend(node, ref, patterns)
	combine(node, ref, patterns, false)

	output(result, patterns)
	close(result)
	return result
}

func starExtend(node *Pattern, ref *Reference, stars map[Char]*Pattern) {
	for idx, mpos := range node.Pos.Iter() {
		k := utils.BitScanLeft64(uint64(mpos))
		if k < 0 {
			continue
		}
		char, next, valid := ref.Next(idx, uint(k))
		for valid {
			pat, exists := stars[char]
			if !exists {
				pat = NewPattern(char, node, false, true)
				stars[char] = pat
			}
			pat.Pos.Add(idx, next)
			char, next, valid = ref.Next(idx, next)
		}
	}
}

func Starex(node *Pattern, ref *Reference) Patterns {
	result := NewPatterns()
	patterns := make(map[Char]*Pattern)
	stars := make(map[Char]*Pattern)
	extend(node, ref, patterns)
	starExtend(node, ref, stars)

	output(result, patterns)
	output(result, stars)

	close(result)
	return result
}

func Regex(node *Pattern, ref *Reference) Patterns {
	result := NewPatterns()
	patterns := make(map[Char]*Pattern)
	stars := make(map[Char]*Pattern)

	extend(node, ref, patterns)
	combine(node, ref, patterns, false)
	starExtend(node, ref, stars)
	combine(node, ref, stars, true)

	output(result, patterns)
	output(result, stars)

	close(result)
	return result
}
