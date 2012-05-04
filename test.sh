#!/bin/bash
export GOPATH=`pwd`
export GOBIN=`pwd`/bin

go install spxs

#time ./bin/spxs --procs=4 --conf=conf/spxs.json inp=data/pres.ref ref=data/pres.bg
#time ./bin/spxs --verbose --procs=4 --conf=conf/spxs.json inp=data/pres.ref ref=data/pres.bg  

CMD="./bin/spxs --verbose --procs=32 --conf=conf/spxs.json inp=data/pres.ref ref=data/pres.bg"

time $CMD

exit

exec 4<>live.log
tail -f live.log &

$CMD 2>&4 | tee actual.log

exec 4>&-
kill %1