package mining

import (
	"testing"
	"fmt"
	"encoding/hex"
	"crypto/sha256"

	"github.com/tclchiam/block_n_go/blockchain/chainhash"
	"github.com/tclchiam/block_n_go/blockchain/entity"
)

func buildTransactionId(newId string) entity.TransactionId {
	decoded, _ := hex.DecodeString(newId)

	var id entity.TransactionId
	copy(id[:], decoded[:sha256.Size])
	return id
}

func TestCalculateBlockHash(t *testing.T) {
	buildChainHash := func(hash string) *chainhash.Hash {
		ret, _ := chainhash.NewHashFromStr(hash)
		return ret
	}

	tests := []struct {
		header *entity.BlockHeader
		nonce  int
		output *chainhash.Hash
	}{
		{
			header: entity.NewBlockHeader(
				0,
				&chainhash.EmptyHash,
				[]*entity.Transaction{
					{
						ID: buildTransactionId("69a101b4ab5c06bf126074a32a6eee3c06b5612f59994a9df280ab5c3603c6b8"),
						Outputs: []*entity.Output{
							{Index: 0, Value: 10, PublicKeyHash: []byte("4c0b332404ac6f5d11c96c0f967398ffd94121ce")},
						},
					},
				},
				1515032127,
			),
			nonce:  18193,
			output: buildChainHash("0000d9bbb68fc04dd2b5f34999a35fdf753abfd260d4541967e87e519696a2eb"),
		},
		{
			header: entity.NewBlockHeader(
				0,
				&chainhash.EmptyHash,
				[]*entity.Transaction{
					{
						ID: buildTransactionId("bbbe0e2f0dd48b427fff9e3ac2105aabb070d2fcea365cb40f8c1e84c0b6ce0b"),
						Outputs: []*entity.Output{
							{Index: 0, Value: 10, PublicKeyHash: []byte("65633924d71fb5244d89afe45aabfaf512cfd148")},
						},
					},
				},
				1515036711,
			),
			nonce:  27764,
			output: buildChainHash("0000e61eeb820f5d29e9a2149adb396f4405963ecc0159f6cec52c8de1fbf672"),
		},
		{
			header: entity.NewBlockHeader(
				2,
				buildChainHash("0000745031d715be942d0fc2731fd0f4b4edd340bad2de76a2fa98368be53419"),
				[]*entity.Transaction{
					{
						ID: buildTransactionId("b0093d332b4c5bbb5f3c4aa2c9ada8632f9efb2489799a74c55168f3487ec256"),
						Outputs: []*entity.Output{
							{Index: 0, Value: 4, PublicKeyHash: []byte("ded5a23a73a574f8465db3c154fc4e7fd75c5bdb")},
							{Index: 1, Value: 3, PublicKeyHash: []byte("52a530c258e53e04116f66d9cae093d0a38950a5")},
						},
					},
					{
						ID: buildTransactionId("6ba28ab31ac33141dcf6def7adf601be3229c4aa148cfa69e7036cc2cedf0aff"),
						Outputs: []*entity.Output{
							{Index: 0, Value: 10, PublicKeyHash: []byte("ded5a23a73a574f8465db3c154fc4e7fd75c5bdb")},
						},
					},
				},
				1515037418,
			),
			nonce:  25634,
			output: buildChainHash("000012f029d52582ed6f179f7b949a0eca1e4f3f7115898de1c15c42ac576f42"),
		},
	}

	for index, testParams := range tests {
		err := calculateBlockHashTestSuite(testParams.header, testParams.nonce, testParams.output, index)
		if err != nil {
			t.Error(err)
		}
	}
}

func calculateBlockHashTestSuite(header *entity.BlockHeader, nonce int, output *chainhash.Hash, index int) error {
	const unexpectedResultStr = "CalculateHash #%d got: %s want: %s"

	result := CalculateHash(header, nonce)

	if !output.IsEqual(&result) {
		return fmt.Errorf(unexpectedResultStr, index, result, output)
	}

	return nil
}
