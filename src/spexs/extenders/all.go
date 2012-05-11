package extenders

var All = [...]Desc{
	{"simple",
		"extends using the alphabet",
		wrap(SimpleExtender)},
	{"group",
		"extends using the alphabet and group symbols",
		wrap(GroupExtender)},
	{"star",
		"extends using the alphabet and star extension",
		wrap(StarExtender)},
	{"regexp",
		"extends using the alphabet, group symbols and stars",
		wrap(GroupStarExtender)},
}