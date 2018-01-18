package entity

import (
	"fmt"
	"strings"
)

type Block struct {
	header       *BlockHeader
	transactions Transactions
}

type Blocks []*Block

func NewBlock(header *BlockHeader, transactions Transactions) *Block {
	return &Block{
		header:       header,
		transactions: transactions,
	}
}

func (block *Block) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("\n============ Block ============"))
	lines = append(lines, fmt.Sprintf("Index: %d", block.Index()))
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

func (block *Block) IsEqual(other *Block) bool {
	if !block.header.IsEqual(other.header) {
		return false
	}
	if !block.Transactions().IsEqual(other.Transactions()) {
		return false
	}

	return true
}

func (block *Block) Index() uint64              { return block.header.Index }
func (block *Block) PreviousHash() *Hash        { return block.header.PreviousHash }
func (block *Block) Timestamp() uint64          { return block.header.Timestamp }
func (block *Block) Header() *BlockHeader       { return block.header }
func (block *Block) Transactions() Transactions { return block.transactions }
func (block *Block) TransactionsHash() *Hash    { return block.header.TransactionsHash }
func (block *Block) Hash() *Hash                { return block.header.Hash }
func (block *Block) Nonce() uint64              { return block.header.Nonce }

func (block *Block) IsGenesisBlock() bool { return block.PreviousHash().IsEmpty() }
