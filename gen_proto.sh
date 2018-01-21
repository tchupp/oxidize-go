#!/bin/bash

main() {
  if ! which protoc &>/dev/null; then
    fail "missing protoc binary (you need to install protoc and put it in \$PATH)"
  fi

  # Entities
  protoc --go_out=. encoding/entities.proto

  # gRPC
  protoc --go_out=Mencoding/entities.proto=github.com/tclchiam/block_n_go/encoding,plugins=grpc:. rpc/*.proto
}

main