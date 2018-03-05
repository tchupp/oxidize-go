package testdata

import (
	"testing"

	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/blockchain/engine/mining/proofofwork"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/blockchain/utxo"
	"github.com/tclchiam/oxidize-go/identity"
	"github.com/tclchiam/oxidize-go/storage/memdb"
)

type BlockchainBuilder struct {
	t *testing.T

	beneficiary     *identity.Identity
	chainRepository entity.ChainRepository
	utxoRepository  utxo.Repository
}

func NewBlockchainBuilder(t *testing.T) *BlockchainBuilder {
	return &BlockchainBuilder{t: t}
}

func (b *BlockchainBuilder) WithBeneficiary(beneficiary *identity.Identity) *BlockchainBuilder {
	b.beneficiary = beneficiary
	return b
}

func (b *BlockchainBuilder) WithChainRepository(repository entity.ChainRepository) *BlockchainBuilder {
	b.chainRepository = repository
	return b
}

func (b *BlockchainBuilder) WithUtxoRepository(repository utxo.Repository) *BlockchainBuilder {
	b.utxoRepository = repository
	return b
}

func (b *BlockchainBuilder) Build() *TestBlockchain {
	if b.beneficiary == nil {
		b.beneficiary = identity.RandomIdentity()
	}
	if b.chainRepository == nil {
		b.chainRepository = memdb.NewChainRepository()
	}
	if b.utxoRepository == nil {
		b.utxoRepository = memdb.NewUtxoRepository()
	}

	bc, err := blockchain.Open(b.chainRepository, b.utxoRepository, proofofwork.NewDefaultMiner(b.beneficiary.Address()))
	if err != nil {
		b.t.Fatalf("error opening blockchain: %s", err)
	}
	return &TestBlockchain{t: b.t, Blockchain: bc}
}
