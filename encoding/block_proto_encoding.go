package encoding

import (
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"fmt"
	"github.com/golang/protobuf/proto"
)

type blockProtoEncoder struct{}

var blockProtoEncoderInstance blockProtoEncoder

func BlockProtoEncoder() entity.BlockEncoder {
	return &blockProtoEncoderInstance
}

func (*blockProtoEncoder) EncodeBlock(block *entity.Block) ([]byte, error) {
	message := ToWireBlock(block)

	data, err := proto.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("serializing block to protobuf: %s", err)
	}

	return data, nil
}

func (*blockProtoEncoder) DecodeBlock(input []byte) (*entity.Block, error) {
	message := &Block{}

	err := proto.Unmarshal(input, message)
	if err != nil {
		return nil, fmt.Errorf("deserializing block from protobuf '%s': %s", input, err)
	}

	return FromWireBlock(message)
}
