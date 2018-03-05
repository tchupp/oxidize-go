package testdata

import (
	"testing"

	"github.com/tclchiam/oxidize-go/account"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/blockchain/testdata"
	"github.com/tclchiam/oxidize-go/identity"
)

type AccountEngineBuilder struct {
	t *testing.T
	*testdata.BlockchainBuilder
}

func NewAccountEngineBuilder(t *testing.T) *AccountEngineBuilder {
	return &AccountEngineBuilder{
		t:                 t,
		BlockchainBuilder: testdata.NewBlockchainBuilder(t),
	}
}

func (b *AccountEngineBuilder) WithBeneficiary(beneficiary *identity.Identity) *AccountEngineBuilder {
	b.BlockchainBuilder.WithBeneficiary(beneficiary)
	return b
}

func (b *AccountEngineBuilder) WithRepository(repository entity.ChainRepository) *AccountEngineBuilder {
	b.BlockchainBuilder.WithChainRepository(repository)
	return b
}

func (b *AccountEngineBuilder) Build() *TestAccountEngine {
	testBlockchain := b.BlockchainBuilder.Build()

	return &TestAccountEngine{
		Engine:         account.NewEngine(testBlockchain),
		TestBlockchain: testBlockchain,
	}
}

type TestAccountEngine struct {
	account.Engine
	*testdata.TestBlockchain
}

func (e *TestAccountEngine) AddBalance(address *identity.Address, balance uint64) *TestAccountEngine {
	e.TestBlockchain.AddBalance(address, balance)
	return e
}
