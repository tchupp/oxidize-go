package entity

import (
	"encoding/json"
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/tclchiam/block_n_go/crypto"
)

type (
	txData struct {
		ID      TransactionId        `json:"id"`
		Inputs  []*txSignedInputData `json:"inputs"`
		Outputs []*txOutputData      `json:"outputs"`
		Secret  []byte               `json:"secret"`
	}

	txSignedInputData struct {
		Reference *txOutputReferenceData `json:"reference"`
		Signature []byte                 `json:"signature"`
		PublicKey []byte                 `json:"public_key"`
	}

	txOutputReferenceData struct {
		ID     TransactionId `json:"id"`
		Output *txOutputData `json:"output"`
	}

	txOutputData struct {
		Index         uint   `json:"index"`
		Value         uint   `json:"value"`
		PublicKeyHash []byte `json:"public_key_hash"`
	}
)

type Encoder func(transaction *Transaction) ([]byte, error)
type Decoder func(data []byte) (*Transaction, error)

func EncodeToJson(transaction *Transaction) ([]byte, error) {
	result, err := json.Marshal(toTxData(transaction))
	if err != nil {
		return nil, fmt.Errorf("serializing transaction to json: %s", err)
	}

	return result, nil
}

func DecodeFromJson(input []byte) (*Transaction, error) {
	var data txData

	err := json.Unmarshal(input, data)
	if err != nil {
		return nil, fmt.Errorf("deserializing transaction from json '%s': %s", input, err)
	}

	return fromTxData(&data)
}

func EncodeToGob(transaction *Transaction) ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(toTxData(transaction))
	if err != nil {
		return nil, fmt.Errorf("serializing transaction to gob: %s", err)
	}

	return result.Bytes(), nil
}

func DecodeFromGob(input []byte) (*Transaction, error) {
	var data txData

	decoder := gob.NewDecoder(bytes.NewReader(input))
	err := decoder.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("deserializing transaction from gob '%s': %s", input, err)
	}

	return fromTxData(&data)
}

func toTxData(transaction *Transaction) *txData {
	var inputs []*txSignedInputData
	for _, input := range transaction.Inputs {
		data := toTxInputData(input)
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

func fromTxData(transaction *txData) (*Transaction, error) {
	var inputs []*SignedInput
	for _, input := range transaction.Inputs {
		data, err := fromTxInputData(input)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, data)
	}

	var outputs []*Output
	for _, output := range transaction.Outputs {
		data := fromTxOutputData(output)
		outputs = append(outputs, data)
	}

	return &Transaction{
		ID:      transaction.ID,
		Inputs:  inputs,
		Outputs: outputs,
	}, nil
}

func toTxInputData(input *SignedInput) *txSignedInputData {
	return &txSignedInputData{
		Reference: toTxOutputReferenceData(input.OutputReference),
		Signature: input.Signature.Serialize(),
		PublicKey: input.PublicKey.Serialize(),
	}
}

func fromTxInputData(input *txSignedInputData) (*SignedInput, error) {
	publicKey, err := crypto.DeserializePublicKey(input.PublicKey)
	if err != nil {
		return nil, err
	}

	signature, err := crypto.DeserializeSignature(input.Signature)
	if err != nil {
		return nil, err
	}

	return &SignedInput{
		OutputReference: fromTxOutputReferenceData(input.Reference),
		PublicKey:       publicKey,
		Signature:       signature,
	}, nil
}

func toTxOutputReferenceData(reference *OutputReference) *txOutputReferenceData {
	return &txOutputReferenceData{
		ID:     reference.ID,
		Output: toTxOutputData(reference.Output),
	}
}

func fromTxOutputReferenceData(reference *txOutputReferenceData) *OutputReference {
	return &OutputReference{
		ID:     reference.ID,
		Output: fromTxOutputData(reference.Output),
	}
}

func toTxOutputData(output *Output) *txOutputData {
	return &txOutputData{
		Index:         output.Index,
		Value:         output.Value,
		PublicKeyHash: output.PublicKeyHash,
	}
}

func fromTxOutputData(output *txOutputData) *Output {
	return &Output{
		Index:         output.Index,
		Value:         output.Value,
		PublicKeyHash: output.PublicKeyHash,
	}
}
