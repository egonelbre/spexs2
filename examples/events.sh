#!/bin/bash

SPEXS=../bin/spexs2

time $SPEXS -stats \
	-verbose=true \
	-conf=events/conf.json \
	inp=events/errors.txt \
	ref=events/random.txt
