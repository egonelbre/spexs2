#!/bin/bash
export GOPATH=`pwd`
export GOBIN=`pwd`/bin

REVS=`git log --oneline | wc -l | sed "s/ //g"`
VERSION=`git describe --tags --long`
TIME=`date +%D\\ %H:%M:%S`
echo -e "package main\n\nconst theVersion = \"$VERSION-rev$REVS\"" > src/spxs/version.go

go install spexs
go install spxs