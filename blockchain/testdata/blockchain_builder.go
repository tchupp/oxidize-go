package testdata

import (
	"testing"

	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/blockchain/engine/mining/proofofwork"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/identity"
	"github.com/tclchiam/oxidize-go/storage/memdb"
)

type BlockchainBuilder struct {
	t *testing.T

	beneficiary *identity.Identity
	repository  entity.ChainRepository
}

func NewBlockchainBuilder(t *testing.T) *BlockchainBuilder {
	return &BlockchainBuilder{t: t}
}

func (b *BlockchainBuilder) WithBeneficiary(beneficiary *identity.Identity) *BlockchainBuilder {
	b.beneficiary = beneficiary
	return b
}

func (b *BlockchainBuilder) WithRepository(repository entity.ChainRepository) *BlockchainBuilder {
	b.repository = repository
	return b
}

func (b *BlockchainBuilder) Build() *TestBlockchain {
	if b.beneficiary == nil {
		b.beneficiary = identity.RandomIdentity()
	}
	if b.repository == nil {
		b.repository = memdb.NewChainRepository()
	}

	bc, err := blockchain.Open(b.repository, proofofwork.NewDefaultMiner(b.beneficiary.Address()))
	if err != nil {
		b.t.Fatalf("error opening blockchain: %s", err)
	}
	return &TestBlockchain{t: b.t, bc: bc}
}
