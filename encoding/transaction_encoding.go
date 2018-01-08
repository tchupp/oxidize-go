package encoding

import (
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/crypto"
)

type (
	txData struct {
		ID      entity.TransactionId `json:"id"`
		Inputs  []*txSignedInputData `json:"inputs"`
		Outputs []*txOutputData      `json:"outputs"`
		Secret  []byte               `json:"secret"`
	}

	txSignedInputData struct {
		Reference *txOutputReferenceData `json:"reference"`
		Signature []byte                 `json:"signature"`
		PublicKey []byte                 `json:"public_key"`
	}

	txUnsignedInputData struct {
		Reference *txOutputReferenceData `json:"reference"`
		PublicKey []byte                 `json:"public_key"`
	}

	txOutputReferenceData struct {
		ID     entity.TransactionId `json:"id"`
		Output *txOutputData        `json:"output"`
	}

	txOutputData struct {
		Index         uint   `json:"index"`
		Value         uint   `json:"value"`
		PublicKeyHash []byte `json:"public_key_hash"`
	}
)

func toTxData(transaction *entity.Transaction) *txData {
	var inputs []*txSignedInputData
	for _, input := range transaction.Inputs {
		data := toTxSignedInputData(input)
		inputs = append(inputs, data)
	}

	var outputs []*txOutputData
	for _, output := range transaction.Outputs {
		data := toTxOutputData(output)
		outputs = append(outputs, data)
	}

	return &txData{
		ID:      transaction.ID,
		Inputs:  inputs,
		Outputs: outputs,
		Secret:  transaction.Secret,
	}
}

func fromTxData(transaction *txData) (*entity.Transaction, error) {
	var inputs []*entity.SignedInput
	for _, input := range transaction.Inputs {
		data, err := fromTxSignedInputData(input)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, data)
	}

	var outputs []*entity.Output
	for _, output := range transaction.Outputs {
		data := fromTxOutputData(output)
		outputs = append(outputs, data)
	}

	return &entity.Transaction{
		ID:      transaction.ID,
		Inputs:  inputs,
		Outputs: outputs,
		Secret:  transaction.Secret,
	}, nil
}

func toTxSignedInputData(input *entity.SignedInput) *txSignedInputData {
	return &txSignedInputData{
		Reference: toTxOutputReferenceData(input.OutputReference),
		Signature: input.Signature.Serialize(),
		PublicKey: input.PublicKey.Serialize(),
	}
}

func fromTxSignedInputData(input *txSignedInputData) (*entity.SignedInput, error) {
	publicKey, err := crypto.DeserializePublicKey(input.PublicKey)
	if err != nil {
		return nil, err
	}

	signature, err := crypto.DeserializeSignature(input.Signature)
	if err != nil {
		return nil, err
	}

	return &entity.SignedInput{
		OutputReference: fromTxOutputReferenceData(input.Reference),
		PublicKey:       publicKey,
		Signature:       signature,
	}, nil
}

func toTxUnsignedInputData(input *entity.UnsignedInput) *txUnsignedInputData {
	return &txUnsignedInputData{
		Reference: toTxOutputReferenceData(input.OutputReference),
		PublicKey: input.PublicKey.Serialize(),
	}
}

func fromTxUnsignedInputData(input *txUnsignedInputData) (*entity.UnsignedInput, error) {
	publicKey, err := crypto.DeserializePublicKey(input.PublicKey)
	if err != nil {
		return nil, err
	}

	return &entity.UnsignedInput{
		OutputReference: fromTxOutputReferenceData(input.Reference),
		PublicKey:       publicKey,
	}, nil
}

func toTxOutputReferenceData(reference *entity.OutputReference) *txOutputReferenceData {
	return &txOutputReferenceData{
		ID:     reference.ID,
		Output: toTxOutputData(reference.Output),
	}
}

func fromTxOutputReferenceData(reference *txOutputReferenceData) *entity.OutputReference {
	return &entity.OutputReference{
		ID:     reference.ID,
		Output: fromTxOutputData(reference.Output),
	}
}

func toTxOutputData(output *entity.Output) *txOutputData {
	return &txOutputData{
		Index:         output.Index,
		Value:         output.Value,
		PublicKeyHash: output.PublicKeyHash,
	}
}

func fromTxOutputData(output *txOutputData) *entity.Output {
	return &entity.Output{
		Index:         output.Index,
		Value:         output.Value,
		PublicKeyHash: output.PublicKeyHash,
	}
}
