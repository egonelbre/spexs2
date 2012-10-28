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
