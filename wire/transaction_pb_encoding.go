package wire

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
)

func EncodeTransaction(transaction *entity.Transaction) ([]byte, error) {
	message := ToWireTransaction(transaction)
	data, err := proto.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("serializing transaction to protobuf: %s", err)
	}

	return data, nil
}

func DecodeTransaction(input []byte) (*entity.Transaction, error) {
	message := &Transaction{}

	err := proto.Unmarshal(input, message)
	if err != nil {
		return nil, fmt.Errorf("deserializing transaction from protobuf '%s': %s", input, err)
	}

	return FromWireTransaction(message)
}

func EncodeOutput(output *entity.Output) ([]byte, error) {
	message := ToWireOutput(output)

	data, err := proto.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("serializing transaction output to protobuf: %s", err)
	}

	return data, nil
}

func DecodeOutput(input []byte) (*entity.Output, error) {
	message := &Output{}

	err := proto.Unmarshal(input, message)
	if err != nil {
		return nil, fmt.Errorf("deserializing transaction output from protobuf '%s': %s", input, err)
	}

	return FromWireOutput(message), nil
}

func EncodeSignedInput(input *entity.SignedInput) ([]byte, error) {
	message := ToWireSignedInput(input)

	data, err := proto.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("serializing transaction signed input to protobuf: %s", err)
	}

	return data, nil
}

func DecodeSignedInput(input []byte) (*entity.SignedInput, error) {
	message := &SignedInput{}

	err := proto.Unmarshal(input, message)
	if err != nil {
		return nil, fmt.Errorf("deserializing transaction signed input from protobuf '%s': %s", input, err)
	}

	return FromWireSignedInput(message)
}

func EncodeUnsignedInput(input *entity.UnsignedInput) ([]byte, error) {
	message := ToWireUnsignedInput(input)

	data, err := proto.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("serializing transaction unsigned input to protobuf: %s", err)
	}

	return data, nil
}

func DecodeUnsignedInput(input []byte) (*entity.UnsignedInput, error) {
	message := &UnsignedInput{}

	err := proto.Unmarshal(input, message)
	if err != nil {
		return nil, fmt.Errorf("deserializing transaction unsigned input from protobuf '%s': %s", input, err)
	}

	return FromWireUnsignedInput(message)
}
