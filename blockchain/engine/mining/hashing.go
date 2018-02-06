package mining

import (
	"bytes"
	"encoding/binary"

	"github.com/tclchiam/oxidize-go/blockchain/entity"
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
		intToHex(input.Timestamp),
		intToHex(nonce),
		intToHex(input.Difficulty),
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

func intToHex(num uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, num)
	return enc
}
