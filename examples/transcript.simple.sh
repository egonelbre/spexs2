#!/bin/bash

SPEXS=../bin/spexs2

time $SPEXS -stats \
	-conf=data/transcript/conf.json \
	inp=data/transcript/transcripts.inp
