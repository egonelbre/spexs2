# spexs2

### Quick introduction

What's **spexs**? It's easier to show what it actually does. To start you need a set of sequences - it can be a sequence of letters or words. For let's imagine we have a data set of events some of them lead to errors and some of them ended normally:

```
# without errors
AIBBFCACCADAHABXCHCG
GBACDBHBDAIBHYDIHAAADAFAHFGGDBFFYFZBFBAGDIDDX
CAGZHGBAXHFIGBAFBIABDYBABBFDBAFGGAAAAHHC
CGDCHHAAAABFBDBCHBBFGICDBGGDGCDFIFADCA
... snip ...
# with errors
ADDDBBCYDFCCHXFDDXBAYDYBHACAZE
DXFDIHBXYDBFGGCBHAYBDHZE
IXBBXHBBACYCFHADHGFDACDHCGYABYBHADZE
AHAFFFGABIXBCAYCBBHBDCDDXZE
... snip ...
```

The events `A`, `B`, `C`, `D`, `F`, `G`, `H`, `I`, `X`, `Y`, `Z` are normal events and `E` is the problematic "error" event. So we need to figure out which event pattern causes the problem?

After we prepare the configuration file (explained later) and plug in the data **spexs** gives the result:

```
pattern      errors     all         ratio   p-value
... snip ...
F.*Y.*Z      272/343    522/5000    7.588   1.639e-174
X.*A.*Z      305/343    811/5000    5.478   1.798e-176
X.*B.*Z      304/343    785/5000    5.641   1.540e-178
D.*Y.*Z      281/343    557/5000    7.347   1.066e-180
C.*Y.*Z      285/343    584/5000    7.107   1.056e-181
B.*Y.*Z      292/343    610/5000    6.971   2.370e-187
A.*Y.*Z      300/343    616/5000    7.092   3.163e-198
X.*Z         343/343    1054/5000   4.740   1.276e-215
Y.*Z         343/343    805/5000    6.204   1.009e-249
X.*Y.*Z      343/343    343/5000    14.537  0
```

Which means that we found an interesting pattern `X.*Y.*Z` (using the regular expression notation).

Remember that **spexs** isn't a magical tool that knows which kinds of patterns you are looking for, so you need to configure the tool appropriately.

### Overview

The easiest way to get started is to download an example that is most similar with your problem. Then adjust the configuration and finally tweak until it fits.

The simplest approach to running the tool is:

```
./spexs2 -conf=conf.json inp=data.txt ref=random.txt
```

We need to give it a configuration file to specify how to find patterns. Let's examine the configuration file:

```
{
    "Dataset": {
        "fore" : { "File" : "$inp$" },
        "back" : { "File" : "$ref$" }
    },
    "Reader" : {
        "Method" : "Delimited",
        "Separator" : "",
        "CountSeparator": ""
    },
    "Extension": {
        "Method": "Star",
        "Groups" : {
            "." : { "elements" : "ABCDFGHIXYZ"}
        },
        "Extendable": {
            "PatLength()" : {"max" : 7},
            "PatGroups()" : {"max" : 3},
            "Matches(fore)" : {"min" : 250},
            "NoStartingGroup()" : {},
            "NoTokens()" : {"Tokens":"E"},
        },
        "Outputtable" : {
            "PatLength()" : {"min" : 3},
            "Hyper(fore, back)" : {"max" : $limit:1e-10$},
            "NoEndingGroup()" : {}
        }
    },
    "Output": {
        "SortBy": ["-Hyper(fore,back)", "Matches(fore)"],
        "Count": 20
    },
    "Printer" : {
        "Method" : "Formatted",
        "ShowHeader": true,
        "Reverse": true,
        "Header": "",
        "Format": "Pat?()\tMatchesPropRatio(fore,back)\tHyper(fore,back)\n"
    }
}
```

First the configuration looks a lot like *json* format, but it can also deal with small user errors, like missing trailing commas. As we can see the configuration is divided into 5 parts: Dataset, Reader, Extension, Output and Printer. Each specifies a particular part of the tool.

#### Dataset

The datasets define where we get our input sequences. There are several ways to specify which files to load:

```
"Dataset": {
    "single" : { "File" : "data" },
    "multiple" : { "Files" : ["a.txt", "b.txt", "c.txt"]},
    "filelist" : { "Filelist" : "listing.txt" },
    "cmdline" : {"File" : "$inp$"}
},
```

The key values in the definition are names for the datasets. These names will be used in the configuration file later.

The *single* just loads the sequences in the file "data". The *multiple* and *filelist* loads sequences from multiple files.

The *cmdline* dataset uses special syntax that is used to mix configuration file and command line parameters. If we run the tool with arguments:

```
./spexs2 -conf=conf.json inp=data.txt
```

Then the `$inp$` in the configuration file will be replaced with `data.txt`. You can also provide default values, for cases where the parameter is missing. For example, if we use `$inp:data$` in the configuration and the `inp` parameter isn't specified it will be replaced with `data` in the configuration file.

#### Reader

Reader specifies how to read in the input. "Method" specifies how to read in the sequences and the rest are arguments to the reader. There are currently to ways of reading the input: "Delimited" and "Counted".

```
"Reader" : {
    "Method" : "Delimited",
    "Separator" : "",
    "Skip": "_,.-'",
},
```

"Separator" specifies how to split the sequences into separate tokens. For example using separator "" will break "AB C" into "A", "B", " ", "C" and using separator " " will break "AB C" into "AB", "C". This separator gives the ability to use letters or letter sequences as our token in the patterns.

"Skip" allows to ignore tokens that are in the input sequence. The "Skip" will be split into tokens in the same way as the sequences, using the "Separator".

There is an alternative way for reading the input:

```
"Reader" : {
    "Method" : "Counted",
    "Separator" : "",
    "CountSeparator" : ":"
},
```

This counted reader expects the input to be in the format:

```
number CountSeparator sequence
```

For example:

```
5:ABCDEFGH
2:CBBDAGHE
1:ABCDAD
```

If we know how many times a sequence occured we can use this to provide a weight for the patterns.

#### Extension

Extension part specifies which patterns will be searched and are nice enough for output.

```
"Extension": {
    "Method": "Star",
    "Groups" : {
        "g1" : { "elements" : "ABC"},
        "g2" : { "elements" : "CDE"},
    },
    "Extendable": {
        "PatLength()" : {"max" : 7},
        "PatGroups()" : {"max" : 3},
        "Matches(fore)" : {"min" : 250},
        "NoStartingGroup()" : {},
        "NoTokens()" : {"Tokens":"E"},
    },
    "Outputtable" : {
        "PatLength()" : {"min" : 3},
        "Hyper(fore, back)" : {"max" : $limit:1e-10$},
        "NoEndingGroup()" : {}
    }
},
```

The "Method" specifies which kinds of patterns it will search. Currently there are 4 different ways for finding patterns: Simple, Group, Star, Regex. 

##### Methods

The "Simple" method will find only patterns that are a simple sequence of tokens.

The "Group" method will find patterns that can contain groups specified in the "Groups" parameter. The group means an alternate possibility for a token. For example, in regular expression syntax it can find patterns `A[BC]D`. The "Groups" allows to specify which groups will be searched for.

The "Star" method will find patterns that can contain gaps in the pattern. For example, in regular expression it would look like `A.*D`.

The "Regex" method will find patterns that can contain either groups or gaps. For example, it could look like `A.*[BC]`.

##### Features

Features are the properties of the patterns. For example, how many times the pattern occurs in a dataset, or how long it is, or the probablity of finding such pattern.

Features allow to define filters (explained in the next section) and they can be used for printing the result. To specify a feature we use the name of the filter and use data set names as argument. For example, to find the number of matching sequences in the data set "alpha" we can specify `Matches(alpha)`.

The full list of features can be seen on the command line with `./spexs2 --details`.

There are important features "Binom" and "Hyper" that allow us to find interesting patterns. To determine whether a pattern happens by chance or is, indeed, characteristic to the data we need to use some statistics. We can calculate how probable some pattern is when we compare the number of matches in the data to the number of matches in some random data. If a pattern is found in the input and random data the same number of times, it probably occured by chance. If a pattern is only found in the input data, then it probably is very interesting.

##### Filters

The filters specify how to limit our search space - or more trivially, how to skip looking at the patterns.

The "Extendable" filter specifies our search space, for example here we can add a limit to the pattern length or skip patterns with low number of occurrences.

The "Outputtable" filter specifies the patterns that we allow to be printed. For example we probably do not want very small patterns.

The feature filter defintion looks like:

```
"Feature(dataset)" : { "min" : 0, "max" : $max:13$ }
```

We can use the special syntax to specify limits from the command-line as well.

#### Output

We can use output to specify how many of the "best" patterns will be remembered and how they will be sorted:

```
"Output": {
    "SortBy": ["-Hyper(fore,back)", "Matches(fore)"],
    "Count": 20
},
```

"SortBy" defines how to decide which are the top results. It takes an array of features as the sorting, if the first feature value is same for two patterns it will use the second etc.

"Count" simply specifies the how many top results to keep.

#### Printer

Printer specifies how to show the results in the command line:

```
"Printer" : {
    "ShowHeader": true,
    "Reverse": true,
    "Header": "",
    "Format": "Pat?()\tMatchesPropRatio(fore,back)\tHyper(fore,back)\n"
}
```

"ShowHeader" specifies whether to add an header to the output. "Reverse" allows to reverse the ordering of the output patterns, this is especially useful when testing the tool, since the best results will be printed to the console last.

"Header" is an optional argument for giving nice looking column headers. If the value is empty the "Format" string will be used instead.

"Format" is how we want to print the feature values to the output.


### Tips

To get started find an example configuration file from the examples that is most similar to your problem.

If you want to use it for large datasets, first try running it on a smaller sample to see whether the limits work nicely.

If the process takes too long try putting in more "Extendable" filters or adding more limits. If the tool doesn't find any results try removing some "Extendable" and/or "Outputtable" filters.

### License

```
spexs2
Copyright (C) 2012  University of Tartu

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
```