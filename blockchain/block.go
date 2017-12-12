package blockchain

import (
	"time"
	"github.com/tclchiam/block_n_go/tx"
	"crypto/sha256"
	"bytes"
)

type Block struct {
	Index        int
	PreviousHash []byte
	Timestamp    int64
	Transactions []*tx.Transaction
	Hash         []byte
	Nonce        int
}

func NewGenesisBlock(address string) *Block {
	transaction := tx.NewGenesisCoinbaseTx(address)

	return NewBlock([]*tx.Transaction{transaction}, []byte(nil), 0)
}

func NewBlock(transactions []*tx.Transaction, previousHash []byte, index int) *Block {
	b := Block{
		Index:        index,
		PreviousHash: previousHash,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
	}
	nonce, hash := CalculateProofOfWork(&b)

	return &Block{
		Index:        b.Index,
		PreviousHash: b.PreviousHash,
		Timestamp:    b.Timestamp,
		Transactions: b.Transactions,
		Hash:         hash,
		Nonce:        nonce,
	}
}

func (block *Block) HashTransactions() []byte {
	var transactionHashes [][]byte

	for _, transaction := range block.Transactions {
		transactionHashes = append(transactionHashes, transaction.ID)
	}

	transactionHash := sha256.Sum256(bytes.Join(transactionHashes, []byte{}))

	return transactionHash[:]
}
