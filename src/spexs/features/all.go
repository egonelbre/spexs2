package features

type CreateFunc interface{}

var All = [...]CreateFunc{
	// simple counting
	Seqs, Matches, Occs,
	// ratios and proportions
	MatchesProp, MatchesRatio, OccsRatio, MatchesPropRatio,
	// hypergeometrics
	Hyper, HyperApprox, HyperDown,
	// pattern length related
	PatLength, PatChars, PatGroups, PatStars,
	// only strings
	Pat, PatRegex,
}

var Help = `
:Pattern:
  Pat?()   : pattern as string
  Regex?() : pattern where groups have been expanded

  PatLength() : pattern length
  PatChars()  : count of simple tokens in pattern
  PatGroups() : count of grouping tokens in pattern
  PatStars()  : count of star tokens in pattern

:Counting:
  Seqs(group)    : total count of sequences
  Matches(group) : count of matching sequences
  Occs(group)    : count of occurences in the sequences

:Ratios:
  MatchesProp(A)     = Matches(A)/Total(A)
  MatchesRatio(A, B) = (Matches(A)+1)/(Matches(B)+1)
  OccsRatio(A, B)    = (Occs(A)+1)/(Occs(B)+1)
  MatchesPropRatio(A, B) = 
    ((Matches(A)+1)/(Total(A)+1))/((Matches(B)+1)/(Total(B)+1))

:Statistics:
  Hyper(fore, back) : hypergeometric p-value
  HyperApprox(fore, back): approx. hypergeometric p-value (~5 sig. digits)
  HyperDown(fore, back) : hypergeometric split down
`
