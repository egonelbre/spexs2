package features

var All = [...]interface{}{
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
