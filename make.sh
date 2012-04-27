#!/bin/bash
export GOPATH=`pwd`
export GOBIN=`pwd`/bin

TIME=`date  +%D\\ %H:%M:%S`
SHA=`git log --pretty=format:'%h' -n 1`
echo -e "package main\n\nconst theVersion = \"v0.2 + $TIME g.$SHA\"" > src/spxs/version.go

go install spexs
go install spxs