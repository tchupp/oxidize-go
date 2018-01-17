package txsigning

import (
	"testing"

	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/blockchain/entity/encoding"
	"github.com/tclchiam/block_n_go/identity"
)

func serializeSignatureTestData(input *entity.UnsignedInput, outputs []*entity.Output, encoder entity.TransactionEncoder) ([]byte, error) {
	var data []byte

	bytes, err := encoder.EncodeUnsignedInput(input)
	if err != nil {
		return nil, err
	}

	data = append(data, bytes...)
	for _, output := range outputs {
		bytes, err := encoder.EncodeOutput(output)
		if err != nil {
			return nil, err
		}
		data = append(data, bytes...)
	}

	return data, nil
}

func TestGenerateInputSignature(t *testing.T) {
	spender := identity.RandomIdentity()

	tests := []struct {
		input   *entity.UnsignedInput
		outputs []*entity.Output
	}{
		{
			input: &entity.UnsignedInput{
				OutputReference: &entity.OutputReference{
					ID:     entity.NewHashOrPanic("b0093d332b4c5bbb5f3c4aa2c9ada8632f9efb2489799a74c55168f3487ec256"),
					Output: &entity.Output{Index: 1, Value: 3, PublicKeyHash: []byte("52a530c258e53e04116f66d9cae093d0a38950a5"),},
				},
				PublicKey: spender.PublicKey(),
			},
			outputs: []*entity.Output{
				{Index: 0, Value: 10, PublicKeyHash: []byte("65633924d71fb5244d89afe45aabfaf512cfd148")},
			},
		},
		{
			input: &entity.UnsignedInput{
				OutputReference: &entity.OutputReference{
					ID:     entity.NewHashOrPanic("caf99368d2abd229d6ff7ec5abdbfdfc7c0b2a2938f23fcb5965a30b4d70ebf8"),
					Output: &entity.Output{Index: 5, Value: 18, PublicKeyHash: []byte("31d6128eb6fbb09e477640ed59252e44c779639f"),},
				},
				PublicKey: spender.PublicKey(),
			},
			outputs: []*entity.Output{
				{Index: 0, Value: 4, PublicKeyHash: []byte("31d6128eb6fbb09e477640ed59252e44c779639f")},
				{Index: 1, Value: 3, PublicKeyHash: []byte("5015668b3d27d024441bd23adc757225d45a4ea1")},
				{Index: 2, Value: 7, PublicKeyHash: []byte("5015668b3d27d024441bd23adc757225d45a4ea1")},
			},
		},
	}

	for index, testParams := range tests {
		const testVerifyFailedStr = "GenerateSignature #%d %s was a bad signature"
		const realVerifyResultMismatchStr = "VerifySignature #%d did not agree with test, got: %s, expected: %s"

		signatureData, err := serializeSignatureTestData(testParams.input, testParams.outputs, encoding.TransactionProtoEncoder())
		if err != nil {
			t.Error(err)
			continue
		}

		signature := GenerateSignature(testParams.input, testParams.outputs, spender, encoding.TransactionProtoEncoder())

		signedInput := entity.NewSignedInput(testParams.input, signature)
		testVerifyResult := spender.PublicKey().Verify(signatureData, signature)
		actualVerifyResult := VerifySignature(signedInput, testParams.outputs, encoding.TransactionProtoEncoder())

		if !testVerifyResult {
			t.Errorf(testVerifyFailedStr, index, signature)
		}
		if actualVerifyResult != testVerifyResult {
			t.Errorf(realVerifyResultMismatchStr, index, actualVerifyResult, testVerifyResult)
		}
	}
}
