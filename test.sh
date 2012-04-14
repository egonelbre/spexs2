#!/bin/bash
export GOPATH=`pwd`
export GOBIN=`pwd`/bin

go install spxs

time ./bin/spxs --conf=conf/spxs.json ref=data/dna.gen
