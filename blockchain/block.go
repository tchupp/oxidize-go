package blockchain

import (
	"time"
)

type Block struct {
	Index        int
	PreviousHash []byte
	Timestamp    int64
	Data         []byte
	Hash         []byte
	Nonce        int
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte(nil), 0)
}

func NewBlock(data string, previousHash []byte, index int) *Block {
	b := Block{
		Index:        index,
		PreviousHash: previousHash,
		Timestamp:    time.Now().Unix(),
		Data:         []byte(data),
	}
	nonce, hash := CalculateProofOfWork(&b)

	return &Block{
		Index:        b.Index,
		PreviousHash: b.PreviousHash,
		Timestamp:    b.Timestamp,
		Data:         b.Data,
		Hash:         hash,
		Nonce:        nonce,
	}
}
