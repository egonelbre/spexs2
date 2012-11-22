1.0.0
-------------------

Most important - major overhaul of the configuration file. Backwards compatibility should last til next major revision (which probably isn't very soon).

Here's a full example: https://github.com/egonelbre/spexs/blob/0b9b063aff50662cc5f09e76454471be02ed2c25/examples/proteins/conf.json

"Dataset" : {
		"name" : {
			"File" : "singlefile"
			"Files" : ["fileA", "fileB", "fileC"],
			"FileList" : "filelist"
		}
	}

Dataset section specifies the files you are working with. You can give each dataset an unique name. It can specify single file, multiple files.

For specifying multiple files you can use "Files" or if you are running it from CLI then using "FileList" is probably easier.

There is always available dataset "@" which contains all files.

	"Reader" : {
		"Method" : "Counted",
		"Separator" : "",
		"CountSeparator" : " "
	}

Reader section specifies how the files should be read:

Method "Delimited" just reads each line and splits to tokens using Separator. Method "Counted" can read in the count for a sequence. The count should be first token on the line and separate by CountSeparator


	"Extension": {
		"Method": "Group",
		"Groups" : {
			"." : { "elements" : "ACDEFGHIKLMNPRQSTVWY"}
		},

		"Extendable": {
			"PatGroups()" : {"max" : 3},
			"PatLength()" : {"max" : 6},
			"Matches(fore)" : {"min" : 20},
			"NoStartingGroup()" : {}
		},
		"Outputtable" : {
			"MatchesRatio(fore, back)" : {"min" : 2},
			"HyperApprox(fore, back)" : {"max" : 1e-3},
			"NoStartingGroup" : {},
			"NoEndingGroup()" : {}
		}
	}

Extension specifies how the algorithm behaves.

For Method there are "Simple", "Group", "Star", "Regexp" - look in --details for more information.

Groups specifies which groups to use for extending things. Elements should be separated by the same token specified in the Reader.

Extendable filters which patterns should be further extended, outputtable filters which patterns are eligible for printing.

There are two types of filters: simple filters and feature filters. For simple filters read the detailed help.

Feature filters take arguments "min", "max" to specify the allowed limits. Additionally Features can take arguments with "()". They specify feature on a particular dataset. For example "Matches(fore)" would count the number of matches in the "fore" dataset.

There are two ways to use a feature : "Feature()" and "Feature?()". Without "?" it just returns the numeric value of the feature, but with "?" it means extended info. For example pattern requires a string to be printed hence we must use "Pat?()" to print the extended info.

All filters have an argument "enabled" for disabling/enabling it from command line.

	"Output": {
		"SortBy": ["-Hyper(fore, back)"],
		"Count": 100
	},

Output specifies how many and in which order output the results. Count specifies how many to output (if it's -1, it will output all). SortBy defines by which features to sort the patterns. "-Feature()" will output the minimal elements and "+Feature()" will output the largest elements.

	"Printer" : {
		"Method" : "Formatted",
		"ShowHeader": true,
		"Reverse": true,
		"Header": "",
		"Format": "Pat?()\tPatLength()\tPatGroups()\tPatChars()\tMatches(fore)\tMatches(back)\tMatchesRatio(fore,back)\tHyper(fore,back)\n"
	}

Printer specifies how to print the result. Header can optionally used to print a custom header. Format specifies what features to print.

If Reverse is true the best result will appear at the bottom (useful for running from command line).

Additionally it's possible to use command line to replace some particular part of the configuration.

	"Extendable" : {
		"NoStartingGroup()" : {"enabled" : "$nostart=true$"}
	}

	> spexs2 -conf=conf.json nostart=false

This would replace the $..$ with false. Or if no value is specified it uses the value after "=".



0.9.4
------

features:

* new query features: "match-hyper-optimal-pvalue", "match-hyper-optimal-seqs", "match-hyper-optimal-matches" - this is calculates minimal hypergeometric pvalue over possible linear splits, you can get how many matches were in the top part of the split with "matches" and how many elements from the top were used in the split "seqs"

0.9.3
------

fixes:

* fixed star extension
* renamed spxs to spexs
* better versioning

0.9.0
------

changes:

* in configuration alphabet.groupName.group is now alphabet.groupName.elements
* removed characters definition from configuration, it will be automatically inferred from the input
* added alphabet.separator for alternative inputs. The lines and group elements will be broken based on that separator. See data/text.json for an alternative example.

fixes : 

* fixed append problem introduced in v0.7:82a18155
