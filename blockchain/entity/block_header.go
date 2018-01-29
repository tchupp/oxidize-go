package entity

import (
	"fmt"
	"strings"
)

type HeaderRepository interface {
	Head() (head *BlockHeader, err error)

	Header(hash *Hash) (*BlockHeader, error)

	SaveHeader(*BlockHeader) error
}

type BlockHeader struct {
	Index            uint64
	PreviousHash     *Hash
	Timestamp        uint64
	TransactionsHash *Hash
	Nonce            uint64
	Hash             *Hash
}

type BlockHeaders []*BlockHeader

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
	lines = append(lines, fmt.Sprintf("PreviousHash: %s", header.PreviousHash))
	lines = append(lines, fmt.Sprintf("Timestamp: %d", header.Timestamp))
	lines = append(lines, fmt.Sprintf("TransactionsHash: %s", header.TransactionsHash))
	lines = append(lines, fmt.Sprintf("Nonce: %d", header.Nonce))
	lines = append(lines, fmt.Sprintf("Hash: %s", header.Hash))

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
	if header.Nonce != other.Nonce {
		return false
	}
	if !header.Hash.IsEqual(other.Hash) {
		return false
	}

	return true
}

func (headers BlockHeaders) Add(header *BlockHeader) BlockHeaders {
	return append(headers, header)
}

func (headers BlockHeaders) Hashes() []*Hash {
	var hashes []*Hash
	for _, header := range headers {
		hashes = append(hashes, header.Hash)
	}
	return hashes
}
