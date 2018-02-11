package encoding

import (
	"github.com/golang/protobuf/proto"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/crypto"
)

func ToWireTransaction(transaction *entity.Transaction) *Transaction {
	var inputs []*SignedInput
	for _, input := range transaction.Inputs {
		data := ToWireSignedInput(input)
		inputs = append(inputs, data)
	}

	var outputs []*Output
	for _, output := range transaction.Outputs {
		data := ToWireOutput(output)
		outputs = append(outputs, data)
	}

	return &Transaction{
		Id:      transaction.ID.Slice(),
		Inputs:  inputs,
		Outputs: outputs,
		Secret:  transaction.Secret,
	}
}

func FromWireTransaction(transaction *Transaction) (*entity.Transaction, error) {
	var inputs []*entity.SignedInput
	for _, input := range transaction.GetInputs() {
		data, err := FromWireSignedInput(input)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, data)
	}

	var outputs []*entity.Output
	for _, output := range transaction.GetOutputs() {
		data := FromWireOutput(output)
		outputs = append(outputs, data)
	}

	id, err := entity.NewHash(transaction.GetId())
	if err != nil {
		return nil, err
	}

	return &entity.Transaction{
		ID:      id,
		Inputs:  inputs,
		Outputs: outputs,
		Secret:  transaction.GetSecret(),
	}, nil
}

func ToWireSignedInput(input *entity.SignedInput) *SignedInput {
	return &SignedInput{
		Reference: ToWireOutputReference(input.OutputReference),
		Signature: input.Signature.Serialize(),
		PublicKey: input.PublicKey.Serialize(),
	}
}

func FromWireSignedInput(input *SignedInput) (*entity.SignedInput, error) {
	publicKey, err := crypto.DeserializePublicKey(input.GetPublicKey())
	if err != nil {
		return nil, err
	}

	signature, err := crypto.DeserializeSignature(input.GetSignature())
	if err != nil {
		return nil, err
	}

	outputReference, err := FromOutputReference(input.GetReference())
	if err != nil {
		return nil, err
	}

	return &entity.SignedInput{
		OutputReference: outputReference,
		PublicKey:       publicKey,
		Signature:       signature,
	}, nil
}

func ToWireUnsignedInput(input *entity.UnsignedInput) *UnsignedInput {
	return &UnsignedInput{
		Reference: ToWireOutputReference(input.OutputReference),
		PublicKey: input.PublicKey.Serialize(),
	}
}

func FromWireUnsignedInput(input *UnsignedInput) (*entity.UnsignedInput, error) {
	publicKey, err := crypto.DeserializePublicKey(input.GetPublicKey())
	if err != nil {
		return nil, err
	}

	outputReference, err := FromOutputReference(input.GetReference())
	if err != nil {
		return nil, err
	}

	return &entity.UnsignedInput{
		OutputReference: outputReference,
		PublicKey:       publicKey,
	}, nil
}

func ToWireOutputReference(reference *entity.OutputReference) *OutputReference {
	return &OutputReference{
		Id:     reference.ID.Slice(),
		Output: ToWireOutput(reference.Output),
	}
}

func FromOutputReference(reference *OutputReference) (*entity.OutputReference, error) {
	id, err := entity.NewHash(reference.GetId())
	if err != nil {
		return nil, err
	}

	return &entity.OutputReference{
		ID:     id,
		Output: FromWireOutput(reference.GetOutput()),
	}, nil
}

func ToWireOutput(output *entity.Output) *Output {
	return &Output{
		Index:         proto.Uint32(output.Index),
		Value:         proto.Uint64(output.Value),
		PublicKeyHash: output.PublicKeyHash,
	}
}

func FromWireOutput(output *Output) *entity.Output {
	return &entity.Output{
		Index:         output.GetIndex(),
		Value:         output.GetValue(),
		PublicKeyHash: output.GetPublicKeyHash(),
	}
}
