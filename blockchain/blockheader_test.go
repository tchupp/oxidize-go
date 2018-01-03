package blockchain

import (
	"bytes"
	"testing"
	"github.com/tclchiam/block_n_go/chainhash"
)

func TestNewGenesisBlockHeader(t *testing.T) {
	const address = "1DtTgiTBGmtLgawnm3ar5FE526FbKMKzCn"
	genesisBlock := NewGenesisBlockHeader(address)

	if bytes.Compare(genesisBlock.PreviousHash.Slice(), chainhash.EmptyHash.Slice()) != 0 {
		t.Fatalf("Genesis block has bad PreviousHash, expected [%s], but was [%s]", chainhash.EmptyHash, genesisBlock.PreviousHash)
	}
	if genesisBlock.Index != 0 {
		t.Fatalf("Genesis block has bad Index, expected %d, but was %d", 0, genesisBlock.Index)
	}
	if len(genesisBlock.Transactions) != 1 {
		t.Fatalf("Genesis block has bad Transactions, %s", genesisBlock.Transactions)
	}
}
