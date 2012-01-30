package spexs

func trieSimpleExtend(n TrieNode, ref Reference, patterns map[Char]Pattern) {
  n.Pos.Iterate( func( pos Pos) {
  	char, next, valid := ref.Next(pos)
    if !valid { continue }
    pat, exists := pats[char]
    if !exists {
      pat = NewTrieNode(char, n)
      patterns[c] = pat
    }
    pat.Pos.Add( next )
  })
}

func SimpleExtender(p Pattern, ref Reference) Patterns {
  result := make(Patterns)
  pats := make( map[Char] Pattern )

  trieSimpleExtend(p.(TrieNode), ref, pats)
  
  for _, pat := pats {
    result <- pat
  }

  return result
}

func trieGroupCombine(n TrieNode, ref Reference, patterns map[Char] Pattern) {
  for g, _ := range ref.Groups(){
    pat := NewTrieNode(g.Id, n)
    patterns[g.Id] = pat
    for char := range g.Chars {
      pats[g.Id].Pos.AddSet( pats[c].Pos )
    }
  }
}

func GroupExtender(p Pattern, ref Reference) Patterns {
  result := make(Patterns)
  pats := make( map[AlphabetChar] Pattern )

  trieSimpleExtend(p.(TrieNode), ref, pats)
  trieGroupCombine(p.(TrieNode), ref, pats)

  for _, pat := pats {
    result <- pat
  }

  return result
}

func trieStarExtend( p TrieNode, ref Reference, stars map[Char] Pattern ) {
  lowest := map[int] Pos

  n.Pos.Iterate( func( p Pos ) {
  	idx, _ := PosDecode(p)
  	cur, exists := lowestPos[idx]
  	if !exists || p < cur {
  		lowestPos[idx] = p
  	}
  })

  for _, p := range lowest {
  	next := p
  	for char, next, valid := ref.Next(p); valid {
  		pat, exists := stars[char]
	    if !exists {
	      pat = NewTrieNode(char, n)
	      pat.IsStar = true
	      stars[c] = pat
	    }
	    pat.Pos.Add( next )
  	}
  }
}

func StarExtender(p Pattern, ref Reference) Patterns {
  result := make(Patterns)
  patterns := make( map[AlphabetChar] Pattern )
  stars := make( map[AlphabetChar] Pattern )

  trieSimpleExtend(p.(TrieNode), ref, patterns)
  trieStarExtend(p.(TrieNode), ref, stars)

  for _, pat := patterns {
    result <- patterns
  }

  return result
}

func GroupStarExtender(p Pattern, ref Reference) Patterns {
  result := make(Patterns)
  patterns := make( map[AlphabetChar] Pattern )
  stars := make( map[AlphabetChar] Pattern )

  trieSimpleExtend(p.(TrieNode), ref, patterns)
  trieGroupCombine(p.(TrieNode), ref, patterns)
  trieStarExtend(p.(TrieNode), ref, stars)
  trieGroupCombine(p.(TrieNode), ref, stars)

  for _, pat := patterns {
    result <- patterns
  }

  return result
}
