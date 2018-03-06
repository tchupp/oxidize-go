#!/bin/bash

main() {
  if ! which protoc &>/dev/null; then
    fail "missing protoc binary (you need to install protoc and put it in \$PATH)"
  fi

  # Entities
  protoc --go_out=. encoding/entities.proto

  # gRPC
  protoc --go_out=Mencoding/entities.proto=github.com/tclchiam/oxidize-go/encoding,plugins=grpc:. blockchain/blockrpc/sync.proto
  protoc --go_out=plugins=grpc:. p2p/discovery.proto

  # Wallet
  protoc -Iwallet -I. --go_out=Mencoding/entities.proto=github.com/tclchiam/oxidize-go/encoding,plugins=grpc:wallet/rpc wallet.proto
}

main