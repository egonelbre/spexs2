#!/bin/bash
export GOPATH=`pwd`

time ./bin/spxs --conf=conf/spxs.json ref=data/dna.gen
