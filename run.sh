#!/bin/sh

TARGET=$1

DIR=$(cd $(dirname $0) && pwd)
cd ${DIR}

GOPATH=${DIR}/${TARGET}
go run ./${TARGET}/src/app/main.go
