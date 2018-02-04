package mining

import (
	"bytes"
	"encoding/binary"
	"math/big"

	"github.com/tclchiam/block_n_go/blockchain/entity"
)

const (
	targetBits = 16
	hashLength = 256
)

var (
	target = big.NewInt(1).Lsh(big.NewInt(1), uint(hashLength-targetBits))
)

type BlockHashingInput struct {
	Index            uint64
	PreviousHash     *entity.Hash
	Timestamp        uint64
	TransactionsHash *entity.Hash
}

func CalculateBlockHash(input *BlockHashingInput, nonce uint64) *entity.Hash {
	rawBlockContents := [][]byte{
		input.PreviousHash.Slice(),
		input.TransactionsHash.Slice(),
		intToHex(input.Timestamp),
		intToHex(nonce),
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

func HashValid(hash *entity.Hash) bool {
	if hash == nil {
		return false
	}

	return new(big.Int).SetBytes(hash.Slice()).Cmp(target) == -1
}

func intToHex(num uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, num)
	return enc
}
