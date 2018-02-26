package proofofwork

import (
	"testing"

	"github.com/tclchiam/oxidize-go/blockchain/engine/mining"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/identity"
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
	parent    = &entity.BlockHeader{
		Index:            0,
		PreviousHash:     &entity.EmptyHash,
		Timestamp:        timestamp,
		TransactionsHash: &entity.EmptyHash,
		Nonce:            9330,
		Hash:             entity.NewHashOrPanic("00008623b2c8806d056cb4ab9a5c3a57d9f36c017aa6c40fed5767249dcd10a8"),
		Difficulty:       4,
	}

	beneficiary = identity.RandomIdentity().Address()
)

func BenchmarkNewMiner_8(b *testing.B) {
	miner := NewMiner(8, beneficiary).(*miner)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		miner.mineBlock(parent, transactions, timestamp)
	}
}

func BenchmarkNewMiner_16(b *testing.B) {
	miner := NewMiner(16, beneficiary).(*miner)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		miner.mineBlock(parent, transactions, timestamp)
	}
}

func BenchmarkNewMiner_32(b *testing.B) {
	miner := NewMiner(32, beneficiary).(*miner)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		miner.mineBlock(parent, transactions, timestamp)
	}
}

func BenchmarkNewMiner_64(b *testing.B) {
	miner := NewMiner(64, beneficiary).(*miner)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		miner.mineBlock(parent, transactions, timestamp)
	}
}

func BenchmarkNewMiner_128(b *testing.B) {
	miner := NewMiner(128, beneficiary).(*miner)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		miner.mineBlock(parent, transactions, timestamp)
	}
}

func BenchmarkNewMiner_256(b *testing.B) {
	miner := NewMiner(256, beneficiary).(*miner)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		miner.mineBlock(parent, transactions, timestamp)
	}
}

func BenchmarkNewDefaultMiner(b *testing.B) {
	miner := NewDefaultMiner(beneficiary).(*miner)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		miner.mineBlock(parent, transactions, timestamp)
	}
}

func TestMiner_MineBlock(t *testing.T) {
	expectedHeader := entity.NewBlockHeader(
		parent.Index+1,
		parent.Hash,
		mining.CalculateTransactionsHash(transactions),
		timestamp,
		12408,
		entity.NewHashOrPanic("0000aff0376b81589ead4920032fd032f0f4257da9ef47af4810ba12d9acc059"),
		parent.Difficulty,
	)
	expected := entity.NewBlock(expectedHeader, transactions)

	miner := NewDefaultMiner(beneficiary).(*miner)
	newBlock := miner.mineBlock(parent, transactions, timestamp)

	if !expected.IsEqual(newBlock) {
		t.Errorf("expected blocks to be equal. got - %s, wanted - %s", newBlock, expected)
	}
}
