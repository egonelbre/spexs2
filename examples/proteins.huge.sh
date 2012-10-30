#!/bin/bash

SPEXS=../bin/spexs2

time $SPEXS -stats \
	-conf=proteins/conf.json \
	inp=proteins/g21.inp \
	ref=proteins/g27.ref
