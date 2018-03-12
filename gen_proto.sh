#!/bin/bash

main() {
  if ! which protoc &>/dev/null; then
    fail "missing protoc binary (you need to install protoc and put it in \$PATH)"
  fi

  # blockchain
  protoc -I_proto -I. --go_out=encoding blockchain_entities.proto
  protoc -I_proto -I. --go_out=M_proto/blockchain_entities.proto=github.com/tclchiam/oxidize-go/encoding,plugins=grpc:blockchain/blockrpc blockchain_service.proto

  # p2p
  protoc -I_proto -I. --go_out=plugins=grpc:p2p node_discovery_service.proto

  # wallet
  protoc -I_proto -I. --go_out=M_proto/blockchain_entities.proto=github.com/tclchiam/oxidize-go/encoding,plugins=grpc:wallet/rpc wallet_service.proto
}

main