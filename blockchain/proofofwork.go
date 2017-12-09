package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"encoding/binary"
	"log"
)

const (
	maxNonce   = math.MaxInt64
	targetBits = 16
	hashLength = 256
)

func CalculateProofOfWork(block *UncommittedBlock) (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	target := big.NewInt(1)
	target.Lsh(target, uint(hashLength-targetBits))

	fmt.Printf("Mining the block containing \"%s\"\n", block.Data)
	for nonce < maxNonce {
		hash = hashBlock(block, nonce)

		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

func (block *CommittedBlock) Validate() bool {
	var hashInt big.Int

	target := big.NewInt(1)
	target.Lsh(target, uint(hashLength-targetBits))

	rawBlockContents := [][]byte{
		block.PreviousHash,
		block.Data,
		intToHex(block.Timestamp),
		intToHex(int64(targetBits)),
		intToHex(int64(block.Nonce)),
	}
	rawBlockData := bytes.Join(rawBlockContents, []byte{})
	hash := sha256.Sum256(rawBlockData)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(target) == -1
}

func hashBlock(block *UncommittedBlock, nonce int) [32]byte {
	rawBlockContents := [][]byte{
		block.PreviousHash,
		block.Data,
		intToHex(block.Timestamp),
		intToHex(int64(targetBits)),
		intToHex(int64(nonce)),
	}
	rawBlockData := bytes.Join(rawBlockContents, []byte{})
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
