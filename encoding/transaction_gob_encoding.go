package encoding

import (
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"bytes"
	"encoding/gob"
	"fmt"
)

type transactionGobEncoder struct{}

var transactionGobEncoderInstance transactionGobEncoder

func NewTransactionGobEncoder() entity.TransactionEncoder {
	return &transactionGobEncoderInstance
}

func (*transactionGobEncoder) EncodeTransaction(transaction *entity.Transaction) ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(toTxData(transaction))
	if err != nil {
		return nil, fmt.Errorf("serializing transaction to gob: %s", err)
	}

	return result.Bytes(), nil
}

func (*transactionGobEncoder) DecodeTransaction(input []byte) (*entity.Transaction, error) {
	var data txData

	decoder := gob.NewDecoder(bytes.NewReader(input))
	err := decoder.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("deserializing transaction from gob '%s': %s", input, err)
	}

	return fromTxData(&data)
}

func (*transactionGobEncoder) EncodeOutput(output *entity.Output) ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(toTxOutputData(output))
	if err != nil {
		return nil, fmt.Errorf("serializing transaction output to gob: %s", err)
	}

	return result.Bytes(), nil
}

func (*transactionGobEncoder) DecodeOutput(input []byte) (*entity.Output, error) {
	var data txOutputData

	decoder := gob.NewDecoder(bytes.NewReader(input))
	err := decoder.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("deserializing transaction output from gob '%s': %s", input, err)
	}

	return fromTxOutputData(&data), nil
}

func (*transactionGobEncoder) EncodeSignedInput(input *entity.SignedInput) ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(toTxSignedInputData(input))
	if err != nil {
		return nil, fmt.Errorf("serializing transaction signed input to gob: %s", err)
	}

	return result.Bytes(), nil
}

func (*transactionGobEncoder) DecodeSignedInput(input []byte) (*entity.SignedInput, error) {
	var data txSignedInputData

	decoder := gob.NewDecoder(bytes.NewReader(input))
	err := decoder.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("deserializing transaction signed input from gob '%s': %s", input, err)
	}

	return fromTxSignedInputData(&data)
}

func (*transactionGobEncoder) EncodeUnsignedInput(input *entity.UnsignedInput) ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(toTxUnsignedInputData(input))
	if err != nil {
		return nil, fmt.Errorf("serializing transaction unsigned input to gob: %s", err)
	}

	return result.Bytes(), nil
}

func (*transactionGobEncoder) DecodeUnsignedInput(input []byte) (*entity.UnsignedInput, error) {
	var data txUnsignedInputData

	decoder := gob.NewDecoder(bytes.NewReader(input))
	err := decoder.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("deserializing transaction unsigned input from gob '%s': %s", input, err)
	}

	return fromTxUnsignedInputData(&data)
}
