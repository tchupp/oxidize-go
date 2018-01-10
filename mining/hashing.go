package mining

import (
	"bytes"
	"encoding/binary"
	"math/big"

	"github.com/tclchiam/block_n_go/blockchain/chainhash"
	"github.com/tclchiam/block_n_go/blockchain/entity"
)

const (
	targetBits = 16
	hashLength = 256
)

var (
	target = big.NewInt(1).Lsh(big.NewInt(1), uint(hashLength-targetBits))
)

func CalculateHash(header *entity.BlockHeader, nonce uint64) chainhash.Hash {
	rawBlockContents := [][]byte{
		header.PreviousHash.Slice(),
		header.TransactionsHash.Slice(),
		intToHex(header.Timestamp),
		intToHex(nonce),
	}
	rawBlockData := bytes.Join(rawBlockContents, []byte(nil))
	return chainhash.CalculateHash(rawBlockData)
}

func Valid(block *entity.Block) bool {
	hash := CalculateHash(block.Header(), block.Nonce())

	return new(big.Int).SetBytes(hash.Slice()).Cmp(target) == -1
}

func HashValid(hash chainhash.Hash) bool {
	return new(big.Int).SetBytes(hash.Slice()).Cmp(target) == -1
}

func intToHex(num uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, num)
	return enc
}
