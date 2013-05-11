#!/bin/bash

SPEXS=../bin/spexs2

time $SPEXS -stats \
	-conf=yeast/conf.json \
	inp=yeast/Yeast_-600_+2_W_cluster_1599945.fa \
	ref=yeast/Yeast_-600_+2_W_random_1000_all.fa
