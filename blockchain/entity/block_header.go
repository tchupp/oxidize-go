package entity

import (
	"fmt"
	"time"
	"strings"

	"github.com/tclchiam/block_n_go/blockchain/chainhash"
	"bytes"
)

type BlockHeader struct {
	Index            uint64
	PreviousHash     *chainhash.Hash
	Timestamp        uint64
	TransactionsHash *chainhash.Hash
}

func NewGenesisBlockHeader(transactions Transactions) *BlockHeader {
	return NewBlockHeaderNow(0, &chainhash.EmptyHash, transactions)
}

func NewBlockHeaderNow(index uint64, previousHash *chainhash.Hash, transactions Transactions) *BlockHeader {
	return NewBlockHeader(index, previousHash, transactions, uint64(time.Now().Unix()))
}

func NewBlockHeader(index uint64, previousHash *chainhash.Hash, transactions Transactions, timestamp uint64) *BlockHeader {
	transactionsHash := hashTransactions(transactions)

	return &BlockHeader{
		Index:            index,
		PreviousHash:     previousHash,
		Timestamp:        timestamp,
		TransactionsHash: &transactionsHash,
	}
}

func (header BlockHeader) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("\n============ BlockHeader ============"))
	lines = append(lines, fmt.Sprintf("Index: %x", header.Index))
	lines = append(lines, fmt.Sprintf("PreviousHash: %x", header.PreviousHash.Slice()))
	lines = append(lines, fmt.Sprintf("Timestamp: %d", header.Timestamp))
	lines = append(lines, fmt.Sprintf("TransactionsHash: %x", header.TransactionsHash))

	return strings.Join(lines, "\n")
}

func (header *BlockHeader) IsEqual(other *BlockHeader) bool {
	if header.Index != other.Index {
		return false
	}
	if !header.PreviousHash.IsEqual(other.PreviousHash) {
		return true
	}
	if header.Timestamp != other.Timestamp {
		return false
	}
	if !header.TransactionsHash.IsEqual(other.TransactionsHash) {
		return false
	}

	return true
}

func hashTransactions(transactions Transactions) chainhash.Hash {
	var transactionHashes [][]byte

	for _, transaction := range transactions {
		transactionHashes = append(transactionHashes, transaction.ID.Slice())
	}

	return chainhash.CalculateHash(bytes.Join(transactionHashes, []byte{}))
}
