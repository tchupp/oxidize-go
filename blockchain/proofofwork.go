package blockchain

import (
	"bytes"
	"crypto/sha256"
	"math"
	"encoding/binary"
	"log"
	"runtime"
	"github.com/tclchiam/block_n_go/tx"
)

const maxNonce = math.MaxInt64

func CalculateProofOfWork(header *BlockHeader) (*blockSolution) {
	workerCount := runtime.NumCPU()

	solutions := make(chan *blockSolution)
	nonces := make(chan int, workerCount)
	defer close(nonces)

	for worker := 0; worker < workerCount; worker++ {
		go miner(header, nonces, solutions)
	}

	for nonce := 0; nonce < maxNonce; nonce++ {
		select {
		case solution := <-solutions:
			return solution
		default:
			nonces <- nonce
		}
	}

	log.Panic(MaxNonceOverflowError)
	return nil
}

func (block *Block) Validate() bool {
	hash := hashBlock(block.Header(), block.Nonce)
	return validateSolution(hash)
}

func intToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func hashTransactions(transactions []*tx.Transaction) []byte {
	var transactionHashes [][]byte

	for _, transaction := range transactions {
		transactionHashes = append(transactionHashes, transaction.ID[:])
	}

	transactionHash := sha256.Sum256(bytes.Join(transactionHashes, []byte{}))

	return transactionHash[:]
}
