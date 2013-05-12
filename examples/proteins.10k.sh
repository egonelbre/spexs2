#!/bin/bash

SPEXS=../bin/spexs2

time $SPEXS -stats \
	-procs=4 \
	-verbose=true \
	-conf=proteins/conf.json \
	inp=proteins/g21_10k.inp \
	ref=proteins/g27_10k.ref
