#!/bin/bash

SPEXS=../bin/spexs2

time $SPEXS -procs=16 \
	-stats \
	-cpuprofile=spxs.prof \
	-conf=data/text/conf.json \
	inp=data/text/text.inp \
	ref=data/text/text.ref
