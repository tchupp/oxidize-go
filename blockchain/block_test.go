package blockchain

import (
	"bytes"
	"crypto/sha256"
	"math/big"
	"testing"
	"github.com/tclchiam/block_n_go/wallet"
)

func TestNewGenesisBlock(t *testing.T) {
	address := wallet.NewWallet().GetAddress()
	genesisBlock := NewGenesisBlock(address)

	if len(genesisBlock.PreviousHash) != 0 {
		t.Fatalf("Genesis block has bad PreviousHash, expected [%s], but was [%s]", []byte{}, genesisBlock.PreviousHash)
	}
	if genesisBlock.Index != 0 {
		t.Fatalf("Genesis block has bad Index, expected %d, but was %d", 0, genesisBlock.Index)
	}
	if len(genesisBlock.Transactions) != 1 {
		t.Fatalf("Genesis block has bad Transactions, %s", genesisBlock.Transactions)
	}

	if !hasValidHash(genesisBlock) {
		t.Fatalf("Genesis block is not valid.")
	}
}

func hasValidHash(block *Block) bool {
	hash := calculateBlockHash(block)

	target := big.NewInt(1)
	target.Lsh(target, uint(hashLength-targetBits))

	var hashInt big.Int
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(target) == -1
}

func calculateBlockHash(block *Block) [32]byte {
	rawBlockContents := [][]byte{
		block.PreviousHash,
		hashTransactions(block),
		intToHex(block.Timestamp),
		intToHex(int64(targetBits)),
		intToHex(int64(block.Nonce)),
	}
	rawBlockData := bytes.Join(rawBlockContents, []byte{})
	return sha256.Sum256(rawBlockData)
}

func hashTransactions(block *Block) []byte {
	var transactionHashes [][]byte

	for _, transaction := range block.Transactions {
		transactionHashes = append(transactionHashes, transaction.ID)
	}

	transactionHash := sha256.Sum256(bytes.Join(transactionHashes, []byte{}))

	return transactionHash[:]
}
