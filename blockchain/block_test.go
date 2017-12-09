package blockchain

import (
	"bytes"
	"crypto/sha256"
	"math/big"
	"testing"
)

func TestNewGenesisBlock(t *testing.T) {
	genesisBlock := NewGenesisBlock()

	if len(genesisBlock.PreviousHash) != 0 {
		t.Fatalf("Genesis block has bad PreviousHash, expected [%s], but was [%s]", []byte{}, genesisBlock.PreviousHash)
	}
	if genesisBlock.Index != 0 {
		t.Fatalf("Genesis block has bad Index, expected %d, but was %d", 0, genesisBlock.Index)
	}
	if string(genesisBlock.Data) != "Genesis Block" {
		t.Fatalf("Genesis block has bad Data, expected \"%s\", but was \"%s\"", "Genesis Block", genesisBlock.Data)
	}

	if !hasValidHash(genesisBlock) {
		t.Fatalf("Genesis block is not valid.")
	}
}

func hasValidHash(block *CommittedBlock) bool {
	hash := calculateBlockHash(block)

	target := big.NewInt(1)
	target.Lsh(target, uint(hashLength-targetBits))

	var hashInt big.Int
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(target) == -1
}

func calculateBlockHash(block *CommittedBlock) [32]byte {
	rawBlockContents := [][]byte{
		block.PreviousHash,
		block.Data,
		intToHex(block.Timestamp),
		intToHex(int64(targetBits)),
		intToHex(int64(block.Nonce)),
	}
	rawBlockData := bytes.Join(rawBlockContents, []byte{})
	return sha256.Sum256(rawBlockData)
}
