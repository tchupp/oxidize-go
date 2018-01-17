package proofofwork

import (
	"testing"

	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/mining"
)

var (
	transactions = entity.Transactions{
		{
			ID: entity.NewHashOrPanic("a34a558abd4599cb63141c556357fba1a777f15fa65835ce190781e2bb2452d9"),
			Outputs: []*entity.Output{
				{Index: 0, Value: 10, PublicKeyHash: []byte("0afd17a7153fc34cfa18b05322d7916dbb5ea24f")},
			},
			Secret: []byte("aa407e4c07e7c2437747ebc07de419351c1737c4bba212481362ecec437f2981"),
		},
	}
	timestamp = uint64(1514479677)
	header    = &entity.BlockHeader{
		Index:            0,
		PreviousHash:     &entity.EmptyHash,
		Timestamp:        timestamp,
		TransactionsHash: mining.CalculateTransactionsHash(transactions),
		Nonce:            9330,
		Hash:             entity.NewHashOrPanic("00008623b2c8806d056cb4ab9a5c3a57d9f36c017aa6c40fed5767249dcd10a8"),
	}
	parent = entity.NewBlock(header, transactions)
)

func BenchmarkNewDefaultMiner(b *testing.B) {
	miner := NewDefaultMiner().(*miner)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		miner.mineBlock(parent, transactions, timestamp)
	}
}

func BenchmarkNewMiner_2(b *testing.B) {
	miner := NewMiner(2).(*miner)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		miner.mineBlock(parent, transactions, timestamp)
	}
}

func BenchmarkNewMiner_4(b *testing.B) {
	miner := NewMiner(4).(*miner)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		miner.mineBlock(parent, transactions, timestamp)
	}
}

func BenchmarkNewMiner_8(b *testing.B) {
	miner := NewMiner(8).(*miner)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		miner.mineBlock(parent, transactions, timestamp)
	}
}

func BenchmarkNewMiner_16(b *testing.B) {
	miner := NewMiner(16).(*miner)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		miner.mineBlock(parent, transactions, timestamp)
	}
}

func BenchmarkNewMiner_32(b *testing.B) {
	miner := NewMiner(32).(*miner)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		miner.mineBlock(parent, transactions, timestamp)
	}
}
