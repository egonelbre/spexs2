#!/bin/bash

SPEXS=../bin/spexs2

time $SPEXS -stats \
	-conf=data/text/conf.json \
	inp=data/text/text.inp \
	ref=data/text/text.ref
