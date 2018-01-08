package entity

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/tclchiam/block_n_go/blockchain/chainhash"
	"github.com/tclchiam/block_n_go/blockchain/tx"
)

type Block struct {
	Index        int
	PreviousHash chainhash.Hash
	Timestamp    int64
	Transactions []*tx.Transaction
	Hash         chainhash.Hash
	Nonce        int
}

func NewBlock(header *BlockHeader, solution *BlockSolution) *Block {
	return &Block{
		Index:        header.Index,
		PreviousHash: header.PreviousHash,
		Timestamp:    header.Timestamp,
		Transactions: header.Transactions,
		Hash:         solution.Hash,
		Nonce:        solution.Nonce,
	}
}

func (block *Block) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("\n============ Block ============"))
	lines = append(lines, fmt.Sprintf("Index: %x", block.Index))
	lines = append(lines, fmt.Sprintf("Hash: %x", block.Hash.Slice()))
	lines = append(lines, fmt.Sprintf("PreviousHash: %x", block.PreviousHash.Slice()))
	lines = append(lines, fmt.Sprintf("Nonce: %d", block.Nonce))
	lines = append(lines, fmt.Sprintf("Timestamp: %d", block.Timestamp))
	lines = append(lines, fmt.Sprintf("Transactions:"))
	block.ForEachTransaction(func(transaction *tx.Transaction) {
		lines = append(lines, transaction.String())
	})

	return strings.Join(lines, "\n")
}

func (block *Block) Header() *BlockHeader {
	return &BlockHeader{
		Transactions: block.Transactions,
		Index:        block.Index,
		PreviousHash: block.PreviousHash,
		Timestamp:    block.Timestamp,
	}
}

func (block *Block) IsGenesisBlock() bool {
	return bytes.Compare(block.PreviousHash.Slice(), chainhash.EmptyHash.Slice()) == 0
}

func (block *Block) ForEachTransaction(consume func(*tx.Transaction)) {
	for _, transaction := range block.Transactions {
		consume(transaction)
	}
}
