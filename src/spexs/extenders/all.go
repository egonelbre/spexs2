package extenders

import . "spexs"

type CreateFunc func(Setup, []byte) Extender

var All = [...]Extender{
	Simple, Group, Star, Regex,
}

func wrap(f Extender) CreateFunc {
	return func(s Setup, data []byte) Extender {
		return f
	}
}

var Help = `
  Simple : uses the sequence tokens to discover the patterns  ( ACCT )
  Group : uses additionally defined groups in Alphabet.Groups ( AC[CT]T )
  Star : uses matching anything in the pattern ( AC.*T )
  Regex : uses both group and star token in the pattern ( A[CT].*T )
`
