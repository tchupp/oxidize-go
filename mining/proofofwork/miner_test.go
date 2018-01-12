package proofofwork_test

import (
	"log"
	"testing"

	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/mining/proofofwork"
)

var (
	transactions = entity.Transactions{
		{
			ID: buildTransactionId("a34a558abd4599cb63141c556357fba1a777f15fa65835ce190781e2bb2452d9"),
			Outputs: []*entity.Output{
				{Index: 0, Value: 10, PublicKeyHash: []byte("0afd17a7153fc34cfa18b05322d7916dbb5ea24f")},
			},
			Secret: []byte("aa407e4c07e7c2437747ebc07de419351c1737c4bba212481362ecec437f2981"),
		},
	}
	blockHeader = entity.NewBlockHeader(0, &entity.EmptyHash, transactions, uint64(1514479677))
)

func buildTransactionId(newId string) *entity.Hash {
	id, err := entity.NewHashFromString(newId)
	if err != nil {
		log.Panic(err)
	}

	return id
}

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
