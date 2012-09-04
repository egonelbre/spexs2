=======================================================
Using Closures to Deal with Algorithmic Configurability
=======================================================

Abstract
--------

Dealing with algorithms that have high configurability is a common problem in bioinformatics.


Example: SPEXS
--------------

SPEXS is an application designed to search patterns in a pattern sequence.

input  : dna sequences
output : 

function spexs( reference, input, output, extender, extendable, outputtable ){
	for pattern := input.get() {
		patterns := extender(pattern, reference)
		for potential := range patterns {
			if extendable(potential) {
				input.put( potential )
			}
			if outputtable(potential) {
				results.put( potential )
			}
		}
	}
}

Although most scientist are capable of formulating this as a generic version, the implementation is lacking.

Most code is riddled with 

	#ifdef DEBUG
		fprintf(blah)
	#elseif
		fprintf(else)
	#endif

or

	if(UseExtenderStar){
		...
	} else if (UseExtenderGroup){
		...
	} else if (UseExtenderSimple){
		...
	}


then

	if(outputFormat == "sqr"){

	} else if (outputFormat == "mgr"){
		
	}