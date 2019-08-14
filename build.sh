#!/usr/bin/env bash


CURDIR=`pwd`
OLDGOPATH="$GOPATH"
export GOPATH="$CURDIR":"$GOPATH"

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
gofmt -w src
go install vc-export
export GOPATH="$OLDGOPATH"

