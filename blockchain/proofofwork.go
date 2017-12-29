package blockchain

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
	"encoding/binary"
	"log"
	"runtime"
)

const (
	maxNonce   = math.MaxInt64
	targetBits = 16
	hashLength = 256
)

var (
	target = big.NewInt(1).Lsh(big.NewInt(1), uint(hashLength-targetBits))
)

func CalculateProofOfWork(block *Block) (*blockSolution) {
	workerCount := runtime.NumCPU()

	solutions := make(chan *blockSolution)
	nonces := make(chan int, workerCount)
	defer close(nonces)

	for worker := 0; worker < workerCount; worker++ {
		go miner(block, nonces, solutions)
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
	var hashInt big.Int

	hash := hashBlock(block, block.Nonce)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(target) == -1
}

func hashBlock(block *Block, nonce int) [32]byte {
	rawBlockContents := [][]byte{
		block.PreviousHash,
		block.HashTransactions(),
		intToHex(block.Timestamp),
		intToHex(int64(targetBits)),
		intToHex(int64(nonce)),
	}
	rawBlockData := bytes.Join(rawBlockContents, []byte(nil))
	return sha256.Sum256(rawBlockData)
}

func intToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
