#! /bin/bash

PATH=$PATH:$(go env GOPATH)/bin
protodir=../../protos
outdir=./internal/api/pb

[[ -e "$outdir" ]] || mkdir -p $outdir

protoc --proto_path=$protodir \
    --go_out=./$outdir --go_opt=paths=source_relative \
    --go-grpc_out=./$outdir --go-grpc_opt=paths=source_relative $protodir/userplan.proto