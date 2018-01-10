package encoding

import (
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/crypto"
	"github.com/golang/protobuf/proto"
)

func toTransactionData(transaction *entity.Transaction) *Transaction {
	var inputs []*SignedInput
	for _, input := range transaction.Inputs {
		data := toSignedInputData(input)
		inputs = append(inputs, data)
	}

	var outputs []*Output
	for _, output := range transaction.Outputs {
		data := toOutputData(output)
		outputs = append(outputs, data)
	}

	return &Transaction{
		Id:      transaction.ID.Slice(),
		Inputs:  inputs,
		Outputs: outputs,
		Secret:  transaction.Secret,
	}
}

func fromTransactionData(transaction *Transaction) (*entity.Transaction, error) {
	var inputs []*entity.SignedInput
	for _, input := range transaction.GetInputs() {
		data, err := fromSignedInputData(input)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, data)
	}

	var outputs []*entity.Output
	for _, output := range transaction.GetOutputs() {
		data := fromOutputData(output)
		outputs = append(outputs, data)
	}

	return &entity.Transaction{
		ID:      entity.TxIdFromBytes(transaction.GetId()),
		Inputs:  inputs,
		Outputs: outputs,
		Secret:  transaction.GetSecret(),
	}, nil
}

func toSignedInputData(input *entity.SignedInput) *SignedInput {
	return &SignedInput{
		Reference: toOutputReferenceData(input.OutputReference),
		Signature: input.Signature.Serialize(),
		PublicKey: input.PublicKey.Serialize(),
	}
}

func fromSignedInputData(input *SignedInput) (*entity.SignedInput, error) {
	publicKey, err := crypto.DeserializePublicKey(input.GetPublicKey())
	if err != nil {
		return nil, err
	}

	signature, err := crypto.DeserializeSignature(input.GetSignature())
	if err != nil {
		return nil, err
	}

	return &entity.SignedInput{
		OutputReference: fromOutputReferenceData(input.GetReference()),
		PublicKey:       publicKey,
		Signature:       signature,
	}, nil
}

func toUnsignedInputData(input *entity.UnsignedInput) *UnsignedInput {
	return &UnsignedInput{
		Reference: toOutputReferenceData(input.OutputReference),
		PublicKey: input.PublicKey.Serialize(),
	}
}

func fromUnsignedInputData(input *UnsignedInput) (*entity.UnsignedInput, error) {
	publicKey, err := crypto.DeserializePublicKey(input.GetPublicKey())
	if err != nil {
		return nil, err
	}

	return &entity.UnsignedInput{
		OutputReference: fromOutputReferenceData(input.GetReference()),
		PublicKey:       publicKey,
	}, nil
}

func toOutputReferenceData(reference *entity.OutputReference) *OutputReference {
	return &OutputReference{
		Id:     reference.ID.Slice(),
		Output: toOutputData(reference.Output),
	}
}

func fromOutputReferenceData(reference *OutputReference) *entity.OutputReference {
	return &entity.OutputReference{
		ID:     entity.TxIdFromBytes(reference.GetId()),
		Output: fromOutputData(reference.GetOutput()),
	}
}

func toOutputData(output *entity.Output) *Output {
	return &Output{
		Index:         proto.Uint32(output.Index),
		Value:         proto.Uint32(output.Value),
		PublicKeyHash: output.PublicKeyHash,
	}
}

func fromOutputData(output *Output) *entity.Output {
	return &entity.Output{
		Index:         output.GetIndex(),
		Value:         output.GetValue(),
		PublicKeyHash: output.GetPublicKeyHash(),
	}
}
