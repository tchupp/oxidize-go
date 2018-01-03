package blockchain

import (
	"testing"

	"github.com/tclchiam/block_n_go/tx"
	"github.com/tclchiam/block_n_go/chainhash"
)

const (
	address   = "1Nar3wkAkLopB2aW1coSDUZcMWkAxP66JK"
	timestamp = int64(1514479677)
)

var (
	transactions = []*tx.Transaction{tx.NewGenesisCoinbaseTx(address)}
	block        = &BlockHeader{
		Index:        0,
		PreviousHash: chainhash.EmptyHash,
		Timestamp:    timestamp,
		Transactions: transactions,
	}
)

func BenchmarkCalculateProofOfWork_GenesisBlock(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		CalculateProofOfWork(block)
	}
}
