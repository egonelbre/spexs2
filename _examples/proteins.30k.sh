#!/bin/bash

go build ../

time spexs2 -stats \
	-procs=1 \
	-verbose=true \
	-conf=proteins/conf.json \
	inp=proteins/g21_30k.inp \
	ref=proteins/g27_30k.ref
