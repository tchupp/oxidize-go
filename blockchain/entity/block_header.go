package entity

import (
	"fmt"
	"strings"
)

type BlockHeader struct {
	Index            uint64
	PreviousHash     *Hash
	Timestamp        uint64
	TransactionsHash *Hash
	Nonce            uint64
	Hash             *Hash
}

func NewBlockHeader(index uint64, previousHash *Hash, transactionsHash *Hash, timestamp uint64, nonce uint64, hash *Hash) *BlockHeader {
	return &BlockHeader{
		Index:            index,
		PreviousHash:     previousHash,
		Timestamp:        timestamp,
		TransactionsHash: transactionsHash,
		Nonce:            nonce,
		Hash:             hash,
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
