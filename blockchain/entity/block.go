package entity

import (
	"fmt"
	"strings"
)

type ChainReader interface {
	BestBlock() (head *Block, err error)
	BlockByHash(hash *Hash) (*Block, error)
	BlockByIndex(index uint64) (*Block, error)

	BestHeader() (head *BlockHeader, err error)
	HeaderByHash(hash *Hash) (*BlockHeader, error)
	HeaderByIndex(index uint64) (*BlockHeader, error)
}

type ChainWriter interface {
	SaveBlock(*Block) error

	SaveHeader(*BlockHeader) error
}

type ChainRepository interface {
	ChainReader
	ChainWriter

	Close() error
}

type Block struct {
	header       *BlockHeader
	transactions Transactions
}

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
	lines = append(lines, fmt.Sprintf("Hash: %s", block.Hash()))
	lines = append(lines, fmt.Sprintf("PreviousHash: %s", block.PreviousHash()))
	lines = append(lines, fmt.Sprintf("TransactionHash: %s", block.header.TransactionsHash))
	lines = append(lines, fmt.Sprintf("Difficulty: %d", block.header.Difficulty))
	lines = append(lines, fmt.Sprintf("Timestamp: %d", block.Timestamp()))
	lines = append(lines, fmt.Sprintf("Nonce: %d", block.Nonce()))
	lines = append(lines, fmt.Sprintf("Transactions:"))
	for _, transaction := range block.Transactions() {
		lines = append(lines, transaction.String())
	}

	return strings.Join(lines, "\n")
}

func (block *Block) IsEqual(other *Block) bool {
	if block == nil && other == nil {
		return true
	}
	if block == nil || other == nil {
		return false
	}
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
