#!/bin/bash
export GOPATH=`pwd`
export GOBIN=`pwd`/bin

go install spxs

#time ./bin/spxs --procs=4 --conf=conf/spxs.json inp=data/pres.ref ref=data/pres.bg
time ./bin/spxs --procs=4 --conf=conf/spxs.json inp=data/pres.ref ref=data/pres.bg
