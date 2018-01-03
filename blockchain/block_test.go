package blockchain

import (
	"bytes"
	"crypto/sha256"
	"math/big"
	"testing"
	"github.com/tclchiam/block_n_go/wallet"
	"github.com/tclchiam/block_n_go/chainhash"
)

func TestNewGenesisBlock(t *testing.T) {
	address := wallet.NewWallet().GetAddress()
	genesisBlock := NewGenesisBlock(address)

	if !genesisBlock.IsGenesisBlock() {
		t.Fatalf("Genesis block does not identify as a genesis block")
	}
	if bytes.Compare(genesisBlock.PreviousHash.Slice(), chainhash.EmptyHash.Slice()) != 0 {
		t.Fatalf("Genesis block has bad PreviousHash, expected [%s], but was [%s]", chainhash.EmptyHash, genesisBlock.PreviousHash)
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
		block.PreviousHash.Slice(),
		hashTransactionsForTest(block),
		intToHex(block.Timestamp),
		intToHex(int64(targetBits)),
		intToHex(int64(block.Nonce)),
	}
	rawBlockData := bytes.Join(rawBlockContents, []byte{})
	return chainhash.CalculateHash(rawBlockData)
}

func hashTransactionsForTest(block *Block) []byte {
	var transactionHashes [][]byte

	for _, transaction := range block.Transactions {
		transactionHashes = append(transactionHashes, transaction.ID[:])
	}

	transactionHash := sha256.Sum256(bytes.Join(transactionHashes, []byte{}))

	return transactionHash[:]
}
