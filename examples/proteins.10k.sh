#!/bin/bash

SPEXS=../bin/spexs2

time $SPEXS -procs=16 \
	-stats \
	-cpuprofile=spxs.prof \
	-conf=data/proteins/conf.json \
	inp=data/proteins/g21_10k.inp \
	ref=data/proteins/g27_10k.ref
