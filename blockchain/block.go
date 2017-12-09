package blockchain

import (
	"time"
	"bytes"
	"encoding/gob"
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

func (block *CommittedBlock) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func DeserializeBlock(d []byte) (*CommittedBlock, error) {
	var block CommittedBlock

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}
