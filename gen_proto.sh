#!/bin/bash

main() {
  if ! which protoc &>/dev/null; then
    fail "missing protoc binary (you need to install protoc and put it in \$PATH)"
  fi

  # Entities
  protoc --go_out=. encoding/entities.proto
}

main