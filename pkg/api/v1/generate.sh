#!/bin/sh

go get github.com/golang/protobuf/protoc-gen-go@v1.3.5
protoc -I. --go_out=plugins=grpc:. *.proto
mv github.com/bhojpur/piro/pkg/api/v1/*.go .
rm -r github.com/bhojpur/piro/pkg/api/v1
