package blockchain

import (
	"time"
)

type UncommittedBlock struct {
	Index        int
	PreviousHash []byte
	Timestamp    int64
	Data         []byte
}

type CommittedBlock struct {
	Index        int
	PreviousHash []byte
	Timestamp    int64
	Data         []byte
	Hash         []byte
	Nonce        int
}

func NewGenesisBlock() *CommittedBlock {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlock(data string, previousHash []byte) *CommittedBlock {
	uncommittedBlock := UncommittedBlock{
		Index:        0,
		PreviousHash: previousHash,
		Timestamp:    time.Now().Unix(),
		Data:         []byte(data),
	}

	return uncommittedBlock.commit()
}

func (b *UncommittedBlock) commit() *CommittedBlock {
	nonce, hash := CalculateProofOfWork(b)

	return &CommittedBlock{
		Index:        b.Index,
		PreviousHash: b.PreviousHash,
		Timestamp:    b.Timestamp,
		Data:         b.Data,
		Hash:         hash,
		Nonce:        nonce,
	}
}
