package encoding

import (
	"fmt"

	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/golang/protobuf/proto"
)

type transactionProtoEncoder struct{}

var transactionProtoEncoderInstance transactionProtoEncoder

func TransactionProtoEncoder() entity.TransactionEncoder {
	return &transactionProtoEncoderInstance
}

func (*transactionProtoEncoder) EncodeTransaction(transaction *entity.Transaction) ([]byte, error) {
	message := toTransactionData(transaction)
	data, err := proto.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("serializing transaction to protobuf: %s", err)
	}

	return data, nil
}

func (*transactionProtoEncoder) DecodeTransaction(input []byte) (*entity.Transaction, error) {
	message := &Transaction{}

	err := proto.Unmarshal(input, message)
	if err != nil {
		return nil, fmt.Errorf("deserializing transaction from protobuf '%s': %s", input, err)
	}

	return fromTransactionData(message)
}

func (*transactionProtoEncoder) EncodeOutput(output *entity.Output) ([]byte, error) {
	message := toOutputData(output)

	data, err := proto.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("serializing transaction output to protobuf: %s", err)
	}

	return data, nil
}

func (*transactionProtoEncoder) DecodeOutput(input []byte) (*entity.Output, error) {
	message := &Output{}

	err := proto.Unmarshal(input, message)
	if err != nil {
		return nil, fmt.Errorf("deserializing transaction output from protobuf '%s': %s", input, err)
	}

	return fromOutputData(message), nil
}

func (*transactionProtoEncoder) EncodeSignedInput(input *entity.SignedInput) ([]byte, error) {
	message := toSignedInputData(input)

	data, err := proto.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("serializing transaction signed input to protobuf: %s", err)
	}

	return data, nil
}

func (*transactionProtoEncoder) DecodeSignedInput(input []byte) (*entity.SignedInput, error) {
	message := &SignedInput{}

	err := proto.Unmarshal(input, message)
	if err != nil {
		return nil, fmt.Errorf("deserializing transaction signed input from protobuf '%s': %s", input, err)
	}

	return fromSignedInputData(message)
}

func (*transactionProtoEncoder) EncodeUnsignedInput(input *entity.UnsignedInput) ([]byte, error) {
	message := toUnsignedInputData(input)

	data, err := proto.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("serializing transaction unsigned input to protobuf: %s", err)
	}

	return data, nil
}

func (*transactionProtoEncoder) DecodeUnsignedInput(input []byte) (*entity.UnsignedInput, error) {
	message := &UnsignedInput{}

	err := proto.Unmarshal(input, message)
	if err != nil {
		return nil, fmt.Errorf("deserializing transaction unsigned input from protobuf '%s': %s", input, err)
	}

	return fromUnsignedInputData(message)
}
