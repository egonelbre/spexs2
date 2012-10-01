#!/bin/bash
export GOPATH=`pwd`
export GOBIN=`pwd`/bin

REVS=`git log --oneline | wc -l | sed "s/ //g"`
VERSION=`git describe --tags --long`
TIME=`date +%D\\ %H:%M:%S`
echo -e "package main\n\nconst (\n\ttheVersion = \"$VERSION-rev$REVS\"\n\ttheBuildTime=\"$TIME\"\n)" > src/spexs2/version.go

go install spexs
go install spexs2