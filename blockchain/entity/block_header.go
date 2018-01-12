package entity

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"
)

type BlockHeader struct {
	Index            uint64
	PreviousHash     *Hash
	Timestamp        uint64
	TransactionsHash *Hash
}

func NewGenesisBlockHeader(transactions Transactions) *BlockHeader {
	return NewBlockHeaderNow(0, &EmptyHash, transactions)
}

func NewBlockHeaderNow(index uint64, previousHash *Hash, transactions Transactions) *BlockHeader {
	return NewBlockHeader(index, previousHash, transactions, uint64(time.Now().Unix()))
}

func NewBlockHeader(index uint64, previousHash *Hash, transactions Transactions, timestamp uint64) *BlockHeader {
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

func hashTransactions(transactions Transactions) Hash {
	var transactionHashes [][]byte

	for _, transaction := range transactions {
		transactionHashes = append(transactionHashes, transaction.ID.Slice())
	}

	return Hash(sha256.Sum256(bytes.Join(transactionHashes, []byte{})))
}
