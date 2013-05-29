#!/bin/bash

SPEXS=../bin/spexs2

time $SPEXS -stats \
	-verbose=true \
	-conf=proteins/conf.json \
	inp=proteins/g21_30k.inp \
	ref=proteins/g27_30k.ref
