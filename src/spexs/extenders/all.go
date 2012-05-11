package extenders

var All = [...]Desc{
	{"simple",
		"extends using the alphabet",
		wrap(Simplex)},
	{"group",
		"extends using the alphabet and group symbols",
		wrap(Groupex)},
	{"star",
		"extends using the alphabet and star extension",
		wrap(Starex)},
	{"regexp",
		"extends using the alphabet, group symbols and stars",
		wrap(Regex)},
}