#!/bin/sh

go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
export PATH="$PATH:$(go env GOPATH)/bin"
protoc -I. --go-grpc_out=. *.proto
mv github.com/bhojpur/piro/pkg/api/v1/*.go .
rm -rf github.com