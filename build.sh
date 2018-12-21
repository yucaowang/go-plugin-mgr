#!/usr/bin/env bash

CURRENT_DIR=`pwd`
OLD_GO_PATH="$GOPATH"  #backup origin gopath
OLD_GO_BIN="$GOBIN"    #backup origin gobin

export GOPATH="$CURRENT_DIR" 
export GOBIN="$CURRENT_DIR/bin"

#code format
gofmt -w src
gofmt -w plugin-src

#build and generate bin
go install src/main.go

# build plugins
go build -o $GOBIN/plugins/kv-json.so.1.0 -buildmode=plugin ./plugin-src/kvstore/json 
go build -o $GOBIN/plugins/kv-memory.so.1.0 -buildmode=plugin ./plugin-src/kvstore/memory 

# copy conf
cp -rf conf $GOBIN/

# recover old environments
export GOPATH="$OLD_GO_PATH"
export GOBIN="$OLD_GO_BIN"
