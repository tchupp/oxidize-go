package proofofwork_test

import (
	"testing"

	"github.com/tclchiam/block_n_go/tx"
	"github.com/tclchiam/block_n_go/chainhash"
	"github.com/tclchiam/block_n_go/blockchain"
	"github.com/tclchiam/block_n_go/proofofwork"
)

const (
	address   = "1Nar3wkAkLopB2aW1coSDUZcMWkAxP66JK"
	timestamp = int64(1514479677)
)

var (
	transactions = []*tx.Transaction{tx.NewGenesisCoinbaseTx(address)}
	blockHeader  = &blockchain.BlockHeader{
		Index:        0,
		PreviousHash: chainhash.EmptyHash,
		Timestamp:    timestamp,
		Transactions: transactions,
	}
)

func BenchmarkNewDefaultMiner(b *testing.B) {
	miner := proofofwork.NewDefaultMiner()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		miner.MineBlock(blockHeader)
	}
}

func BenchmarkNewMiner_2(b *testing.B) {
	miner := proofofwork.NewMiner(2)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		miner.MineBlock(blockHeader)
	}
}

func BenchmarkNewMiner_4(b *testing.B) {
	miner := proofofwork.NewMiner(4)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		miner.MineBlock(blockHeader)
	}
}

func BenchmarkNewMiner_8(b *testing.B) {
	miner := proofofwork.NewMiner(8)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		miner.MineBlock(blockHeader)
	}
}

func BenchmarkNewMiner_16(b *testing.B) {
	miner := proofofwork.NewMiner(16)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		miner.MineBlock(blockHeader)
	}
}

func BenchmarkNewMiner_32(b *testing.B) {
	miner := proofofwork.NewMiner(32)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		miner.MineBlock(blockHeader)
	}
}
