#!/bin/bash

SPEXS=../spexs2

time $SPEXS -stats \
	-verbose=true \
	-conf=dna/small.json \
	inp=dna/dna.small
