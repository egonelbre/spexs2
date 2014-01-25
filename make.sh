#!/bin/bash
export GOPATH=`pwd`
export GOBIN=`pwd`/bin

REVS=`git log --oneline | wc -l | sed "s/ //g"`
VERSION=`git describe --tags --long`
TIME=`date +%D\\ %H:%M:%S`
echo -e "package main\n\nfunc init(){\n\ttheVersion = \"$VERSION-rev$REVS\"\n\ttheBuildTime=\"$TIME\"\n}\n" > src/spexs2/autoVersion.go

go install -gcflags -B spexs
go install -gcflags -B spexs2