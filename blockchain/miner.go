package blockchain

import (
	"math/big"
	"bytes"
	"github.com/tclchiam/block_n_go/chainhash"
)

type blockSolution struct {
	nonce int
	hash  chainhash.Hash
}

const (
	targetBits = 16
	hashLength = 256
)

var (
	target = big.NewInt(1).Lsh(big.NewInt(1), uint(hashLength-targetBits))
)

func miner(header *BlockHeader, nonces <-chan int, solutions chan<- *blockSolution) {
	for nonce := range nonces {
		hash := hashBlock(header, nonce)

		if validateSolution(hash) {
			solutions <- &blockSolution{nonce, hash}
		}
	}
}

func validateSolution(hash chainhash.Hash) bool {
	return new(big.Int).SetBytes(hash[:]).Cmp(target) == -1
}

func hashBlock(header *BlockHeader, nonce int) chainhash.Hash {
	rawBlockContents := [][]byte{
		header.PreviousHash[:],
		hashTransactions(header.Transactions),
		intToHex(header.Timestamp),
		intToHex(int64(targetBits)),
		intToHex(int64(nonce)),
	}
	rawBlockData := bytes.Join(rawBlockContents, []byte(nil))
	return chainhash.CalculateHash(rawBlockData)
}
