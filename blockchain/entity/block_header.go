package entity

import (
	"fmt"
	"time"
	"strings"

	"github.com/tclchiam/block_n_go/blockchain/chainhash"
)

type BlockHeader struct {
	Index        int
	PreviousHash chainhash.Hash
	Timestamp    int64
	Transactions []*Transaction
}

func NewGenesisBlockHeader(address string, encoder TransactionEncoder) *BlockHeader {
	transaction := NewGenesisCoinbaseTx(address, encoder)

	return NewBlockHeader(0, chainhash.EmptyHash, []*Transaction{transaction})
}

func NewBlockHeader(index int, previousHash chainhash.Hash, transactions []*Transaction) *BlockHeader {
	return &BlockHeader{
		Index:        index,
		PreviousHash: previousHash,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
	}
}

func (header BlockHeader) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("\n============ BlockHeader ============"))
	lines = append(lines, fmt.Sprintf("Index: %x", header.Index))
	lines = append(lines, fmt.Sprintf("PreviousHash: %x", header.PreviousHash.Slice()))
	lines = append(lines, fmt.Sprintf("Timestamp: %d", header.Timestamp))
	lines = append(lines, fmt.Sprintf("Transactions:"))
	for _, transaction := range header.Transactions {
		lines = append(lines, transaction.String())
	}

	return strings.Join(lines, "\n")
}
