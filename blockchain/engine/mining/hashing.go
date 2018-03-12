package mining

import (
	"bytes"

	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/conv"
)

type BlockHashingInput struct {
	Index            uint64
	PreviousHash     *entity.Hash
	Timestamp        uint64
	TransactionsHash *entity.Hash
	Difficulty       uint64
}

func CalculateBlockHash(input *BlockHashingInput, nonce uint64) *entity.Hash {
	rawBlockContents := [][]byte{
		input.PreviousHash.Slice(),
		input.TransactionsHash.Slice(),
		conv.U64ToBytes(input.Timestamp),
		conv.U64ToBytes(nonce),
		conv.U64ToBytes(input.Difficulty),
	}
	rawBlockData := bytes.Join(rawBlockContents, []byte(nil))
	hash := calculateHash(rawBlockData)
	return &hash
}

func CalculateHeaderHash(header *entity.BlockHeader) *entity.Hash {
	input := &BlockHashingInput{
		Index:            header.Index,
		PreviousHash:     header.PreviousHash,
		Timestamp:        header.Timestamp,
		TransactionsHash: header.TransactionsHash,
		Difficulty:       header.Difficulty,
	}
	return CalculateBlockHash(input, header.Nonce)
}

func CalculateTransactionsHash(transactions entity.Transactions) *entity.Hash {
	var transactionHashes [][]byte

	for _, transaction := range transactions {
		transactionHashes = append(transactionHashes, transaction.ID.Slice())
	}

	rawTransactionData := bytes.Join(transactionHashes, []byte{})
	hash := calculateHash(rawTransactionData)
	return &hash
}
