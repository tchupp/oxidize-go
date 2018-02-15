package entity

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

type BlockHeader struct {
	Index            uint64
	PreviousHash     *Hash
	Timestamp        uint64
	TransactionsHash *Hash
	Nonce            uint64
	Hash             *Hash
	Difficulty       uint64
}

var GenesisParentHeader = BlockHeader{Index: math.MaxUint64, Hash: &EmptyHash, Difficulty: genesisDifficulty}

func NewBlockHeader(index uint64, previousHash *Hash, transactionsHash *Hash, timestamp uint64, nonce uint64, hash *Hash, difficulty uint64) *BlockHeader {
	return &BlockHeader{
		Index:            index,
		PreviousHash:     previousHash,
		Timestamp:        timestamp,
		TransactionsHash: transactionsHash,
		Nonce:            nonce,
		Hash:             hash,
		Difficulty:       difficulty,
	}
}

func (header BlockHeader) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("\n============ BlockHeader ============"))
	lines = append(lines, fmt.Sprintf("Index: %d", header.Index))
	lines = append(lines, fmt.Sprintf("PreviousHash: %s", header.PreviousHash))
	lines = append(lines, fmt.Sprintf("Timestamp: %d", header.Timestamp))
	lines = append(lines, fmt.Sprintf("TransactionsHash: %s", header.TransactionsHash))
	lines = append(lines, fmt.Sprintf("Difficulty: %d", header.Difficulty))
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

func (header *BlockHeader) IsGenesisBlock() bool { return header.PreviousHash.IsEmpty() }

type BlockHeaders []*BlockHeader

func NewBlockHeaders() BlockHeaders {
	return make(BlockHeaders, 0)
}

func (headers BlockHeaders) Len() int                             { return len(headers) }
func (headers BlockHeaders) Swap(i, j int)                        { headers[i], headers[j] = headers[j], headers[i] }
func (headers BlockHeaders) Less(i, j int) bool                   { return headers[i].Index < headers[j].Index }
func (headers BlockHeaders) Add(header *BlockHeader) BlockHeaders { return append(headers, header) }

func (headers BlockHeaders) Sort() BlockHeaders {
	copies := append(BlockHeaders(nil), headers...)
	sort.Sort(copies)
	return copies
}

func (headers BlockHeaders) Hashes() []*Hash {
	var hashes []*Hash
	for _, header := range headers {
		hashes = append(hashes, header.Hash)
	}
	return hashes
}

func (headers BlockHeaders) IsEqual(other BlockHeaders) bool {
	if headers == nil && other == nil {
		return true
	}
	if headers == nil || other == nil {
		return false
	}
	if headers.Len() == 0 && other.Len() == 0 {
		return true
	}
	if headers.Len() != other.Len() {
		return false
	}

	other = other.Sort()
	self := headers.Sort()
	for i := 0; i < self.Len(); i++ {
		if !self[i].Hash.IsEqual(other[i].Hash) {
			return false
		}
	}

	return true
}
