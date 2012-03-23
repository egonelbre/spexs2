#!/bin/bash
export GOPATH=`pwd`
go build spxs
time spxs --ref=dna.gen --chars=dna.set --extender=regexp --fitness=complexity --limit=4000 --procs=8