package blockchain

import (
	"github.com/tclchiam/block_n_go/tx"
	"github.com/tclchiam/block_n_go/chainhash"
	"time"
)

type BlockHeader struct {
	Index        int
	PreviousHash chainhash.Hash
	Timestamp    int64
	Transactions []*tx.Transaction
}

func NewGenesisBlockHeader(address string) *BlockHeader {
	transaction := tx.NewGenesisCoinbaseTx(address)

	return newBlockHeader(0, chainhash.EmptyHash, []*tx.Transaction{transaction})
}

func newBlockHeader(index int, previousHash chainhash.Hash, transactions []*tx.Transaction) *BlockHeader {
	return &BlockHeader{
		Index:        index,
		PreviousHash: previousHash,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
	}
}
