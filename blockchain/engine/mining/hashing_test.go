package mining

import (
	"fmt"
	"testing"

	"github.com/tclchiam/oxidize-go/blockchain/entity"
)

const unexpectedResultStr = "calculateHash #%d got: %s, want: %s"

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
				Difficulty:       4,
			},
			nonce:  221015,
			output: entity.NewHashOrPanic("00001fb17ca0a8622aead1110c57ef46baf4835bf6b2c385bb2389aa1a6ba11b"),
		},
		{
			input: &BlockHashingInput{
				Index:            0,
				PreviousHash:     &entity.EmptyHash,
				Timestamp:        1515036711,
				TransactionsHash: entity.NewHashOrPanic("bbbe0e2f0dd48b427fff9e3ac2105aabb070d2fcea365cb40f8c1e84c0b6ce0b"),
				Difficulty:       4,
			},
			nonce:  93484,
			output: entity.NewHashOrPanic("00001a4358540913d15cb879a613177f9ca6db6ac846bd81ac80eb7149f37fb0"),
		},
		{
			input: &BlockHashingInput{
				Index:            2,
				PreviousHash:     entity.NewHashOrPanic("0000745031d715be942d0fc2731fd0f4b4edd340bad2de76a2fa98368be53419"),
				Timestamp:        1515037418,
				TransactionsHash: entity.NewHashOrPanic("b0093d332b4c5bbb5f3c4aa2c9ada8632f9efb2489799a74c55168f3487ec256"),
				Difficulty:       4,
			},
			nonce:  111743,
			output: entity.NewHashOrPanic("0000be5fa65c064d12987674d325554d7a26c15381c5381b2998d14d849fb3ef"),
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
	result := CalculateBlockHash(input, nonce)

	if !output.IsEqual(result) {
		return fmt.Errorf(unexpectedResultStr, index, result, output)
	}

	return nil
}
