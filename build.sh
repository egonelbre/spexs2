#!/bin/bash
REVS=`git log --oneline | wc -l | sed "s/ //g"`
VERSION=`git describe --tags --long`
TIME=`date +%D\\ %H:%M:%S`

go build -gcflags "-B" -ldflags "-X main.buildversion \"$VERSION-rev$REVS\" -X main.buildtime \"$TIME\"" .
