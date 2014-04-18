#!/bin/bash

SPEXS=../spexs2

time $SPEXS -stats \
	-conf=text/conf.json \
	inp=text/text.inp \
	ref=text/text.ref
