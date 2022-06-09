#!/usr/bin/bash
set -e

# TODO: only conditionally generate js, if the grpc protoc thingy is
# installed.

PATH=$PATH:~/go/bin
exec protoc \
    --proto_path=proto \
    --go_out=pkg/proto \
    --go-grpc_out=pkg/proto \
    --go-grpc_opt=paths=source_relative \
    --go_opt=paths=source_relative \
    --js_out=js \
    --grpc-web_out=import_style=commonjs,mode=grpcwebtext:js/ \
    proto/tle.proto 
