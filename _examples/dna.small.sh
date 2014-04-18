#!/bin/bash

SPEXS=../spexs2

time $SPEXS -stats \
	-conf=dna/small.json \
	inp=dna/dna.small
