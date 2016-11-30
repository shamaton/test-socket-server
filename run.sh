#!/bin/sh

DIR=$(cd $(dirname $0) && pwd)
cd ${DIR}

GOPATH=${DIR}
go run ./src/app/main.go
