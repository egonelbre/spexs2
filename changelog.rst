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
