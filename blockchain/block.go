package blockchain

import (
	"time"

	"github.com/tclchiam/block_n_go/tx"
	"github.com/tclchiam/block_n_go/chainhash"
	"fmt"
	"strconv"
	"strings"
	"bytes"
)

type Block struct {
	Index        int
	PreviousHash chainhash.Hash
	Timestamp    int64
	Transactions []*tx.Transaction
	Hash         chainhash.Hash
	Nonce        int
}

type BlockHeader struct {
	Index        int
	PreviousHash chainhash.Hash
	Timestamp    int64
	Transactions []*tx.Transaction
}

func NewGenesisBlock(address string) *Block {
	transaction := tx.NewGenesisCoinbaseTx(address)

	return NewBlock([]*tx.Transaction{transaction}, chainhash.EmptyHash, 0)
}

func NewBlock(transactions []*tx.Transaction, previousHash chainhash.Hash, index int) *Block {
	header := &BlockHeader{
		Index:        index,
		PreviousHash: previousHash,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
	}
	solution := CalculateProofOfWork(header)

	return &Block{
		Index:        header.Index,
		PreviousHash: header.PreviousHash,
		Timestamp:    header.Timestamp,
		Transactions: header.Transactions,
		Hash:         solution.hash,
		Nonce:        solution.nonce,
	}
}

func (block *Block) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("\n============ Block ============"))
	lines = append(lines, fmt.Sprintf("Index: %x", block.Index))
	lines = append(lines, fmt.Sprintf("Hash: %x", block.Hash.Slice()))
	lines = append(lines, fmt.Sprintf("PreviousHash: %x", block.PreviousHash.Slice()))
	lines = append(lines, fmt.Sprintf("Nonce: %d", block.Nonce))
	lines = append(lines, fmt.Sprintf("Is valid: %s", strconv.FormatBool(block.Validate())))
	lines = append(lines, fmt.Sprintf("Transactions:"))
	block.ForEachTransaction(func(transaction *tx.Transaction) {
		lines = append(lines, fmt.Sprintf(transaction.String()))
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

func (block *Block) ForEachTransaction(consume func(*tx.Transaction)) error {
	for _, transaction := range block.Transactions {
		consume(transaction)
	}
	return nil
}
