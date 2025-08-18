#!/bin/bash
set -e

PROTO_SRC=./protos
GO_OUT=./protos/gen/go
OPENAPI_OUT=./docs/openapi

[[ -e $GO_OUT ]]      || mkdir -p $GO_OUT
[[ -e $OPENAPI_OUT ]] || mkdir -p $OPENAPI_OUT

protoc -I=$PROTO_SRC \
  --go_out=$GO_OUT --go_opt=paths=source_relative \
  --go-grpc_out=$GO_OUT --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=$GO_OUT --grpc-gateway_opt=paths=source_relative \
  --openapiv2_out=$OPENAPI_OUT \
  $(find $PROTO_SRC -name "*.proto")
