#!/bin/bash

SPEXS=../bin/spexs2

time $SPEXS -procs=16 \
	-stats \
	-cpuprofile=spxs.prof \
	-conf=data/transcript/conf.json \
	inp=data/transcript/transcripts.inp
