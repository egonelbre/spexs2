#!/bin/bash
export GOPATH=`pwd`
export GOBIN=`pwd`/bin

go install spxs

time ./bin/spxs -procs=16 -cpuprofile=spxs.prof -conf=data/g.json inp=data/g21_10k.inp ref=data/g27_10k.ref
