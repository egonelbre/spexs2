#!/bin/bash
export GOPATH=`pwd`
go build spxs
time spxs ref=data/dna.gen
