package mining

import (
	"fmt"
	"testing"

	"github.com/tclchiam/block_n_go/blockchain/entity"
)

func TestCalculateBlockHash(t *testing.T) {
	tests := []struct {
		input  *BlockHashingInput
		nonce  uint64
		output *entity.Hash
	}{
		{
			input: &BlockHashingInput{
				Index:            0,
				PreviousHash:     &entity.EmptyHash,
				Timestamp:        1515032127,
				TransactionsHash: entity.NewHashOrPanic("69a101b4ab5c06bf126074a32a6eee3c06b5612f59994a9df280ab5c3603c6b8"),
			},
			nonce:  59930,
			output: entity.NewHashOrPanic("00005253440ff32184f1793467d02bd7fe57034ddf40eaba597b33221acd9a11"),
		},
		{
			input: &BlockHashingInput{
				Index:            0,
				PreviousHash:     &entity.EmptyHash,
				Timestamp:        1515036711,
				TransactionsHash: entity.NewHashOrPanic("bbbe0e2f0dd48b427fff9e3ac2105aabb070d2fcea365cb40f8c1e84c0b6ce0b"),
			},
			nonce:  121517,
			output: entity.NewHashOrPanic("0000a95c7a2984d5a76c6defd1dfe958786621421fd7ce5ac17923443c54abc0"),
		},
		{
			input: &BlockHashingInput{
				Index:            2,
				PreviousHash:     entity.NewHashOrPanic("0000745031d715be942d0fc2731fd0f4b4edd340bad2de76a2fa98368be53419"),
				Timestamp:        1515037418,
				TransactionsHash: entity.NewHashOrPanic("b0093d332b4c5bbb5f3c4aa2c9ada8632f9efb2489799a74c55168f3487ec256"),
			},
			nonce:  194082,
			output: entity.NewHashOrPanic("00005d057c62d8ec612b8f372e0664a4de34736d66be013b787c920a62aa0ddc"),
		},
	}

	for index, testParams := range tests {
		err := calculateBlockHashTestSuite(testParams.input, testParams.nonce, testParams.output, index)
		if err != nil {
			t.Error(err)
		}
	}
}

func calculateBlockHashTestSuite(input *BlockHashingInput, nonce uint64, output *entity.Hash, index int) error {
	const unexpectedResultStr = "calculateHash #%d got: %s want: %s"

	result := CalculateBlockHash(input, nonce)

	if !output.IsEqual(result) {
		return fmt.Errorf(unexpectedResultStr, index, result, output)
	}

	return nil
}
