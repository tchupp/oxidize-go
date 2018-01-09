package entity

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/tclchiam/block_n_go/blockchain/chainhash"
)

type Block struct {
	header       *BlockHeader
	transactions Transactions
	hash         *chainhash.Hash
	nonce        int
}

func NewBlock(header *BlockHeader, solution *BlockSolution, transactions Transactions) *Block {
	return &Block{
		header:       header,
		transactions: transactions,
		hash:         solution.Hash,
		nonce:        solution.Nonce,
	}
}

func (block *Block) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("\n============ Block ============"))
	lines = append(lines, fmt.Sprintf("Index: %x", block.Index))
	lines = append(lines, fmt.Sprintf("Hash: %x", block.Hash().Slice()))
	lines = append(lines, fmt.Sprintf("PreviousHash: %x", block.PreviousHash().Slice()))
	lines = append(lines, fmt.Sprintf("Timestamp: %d", block.Timestamp()))
	lines = append(lines, fmt.Sprintf("Nonce: %d", block.Nonce()))
	lines = append(lines, fmt.Sprintf("Transactions:"))
	for _, transaction := range block.Transactions() {
		lines = append(lines, transaction.String())
	}

	return strings.Join(lines, "\n")
}

func (block *Block) IsGenesisBlock() bool {
	return bytes.Compare(block.PreviousHash().Slice(), chainhash.EmptyHash.Slice()) == 0
}

func (block *Block) IsEqual(other *Block) bool {
	if !block.header.IsEqual(other.header) {
		return false
	}
	if !block.Hash().IsEqual(other.Hash()) {
		return false
	}
	if block.Nonce() != other.Nonce() {
		return false
	}

	return true
}

func (block *Block) Index() int                    { return block.header.Index }
func (block *Block) PreviousHash() *chainhash.Hash { return block.header.PreviousHash }
func (block *Block) Timestamp() int64              { return block.header.Timestamp }
func (block *Block) Header() *BlockHeader          { return block.header }
func (block *Block) Transactions() Transactions    { return block.transactions }
func (block *Block) Hash() *chainhash.Hash         { return block.hash }
func (block *Block) Nonce() int                    { return block.nonce }
