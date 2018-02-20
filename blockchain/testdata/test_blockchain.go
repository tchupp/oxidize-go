package testdata

import (
	"testing"

	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/encoding"
	"github.com/tclchiam/oxidize-go/identity"
)

type TestBlockchain struct {
	t *testing.T

	bc blockchain.Blockchain
}

func (b *TestBlockchain) AddBalance(address *identity.Address, balance uint64) *TestBlockchain {
	outputs := []*entity.Output{entity.NewOutput(balance, address)}
	tx := entity.NewTx(nil, outputs, encoding.TransactionProtoEncoder())

	block, err := b.bc.MineBlock(entity.Transactions{tx})
	if err != nil {
		b.t.Fatalf("error adding balance: %s", err)
	}

	if err := b.bc.SaveBlock(block); err != nil {
		b.t.Fatalf("error adding balance: %s", err)
	}

	return b
}

func (b *TestBlockchain) ToBlockchain() blockchain.Blockchain {
	return b.bc
}
