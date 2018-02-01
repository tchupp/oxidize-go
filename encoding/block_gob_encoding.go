package encoding

import (
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"bytes"
	"encoding/gob"
	"fmt"
)

type blockGobEncoder struct{}

var blockGobEncoderInstance blockGobEncoder

func BlockGobEncoder() entity.BlockEncoder {
	return &blockGobEncoderInstance
}

func (*blockGobEncoder) EncodeBlock(block *entity.Block) ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(ToWireBlock(block))
	if err != nil {
		return nil, fmt.Errorf("serializing block to gob: %s", err)
	}

	return result.Bytes(), nil
}

func (*blockGobEncoder) DecodeBlock(input []byte) (*entity.Block, error) {
	var data Block

	decoder := gob.NewDecoder(bytes.NewReader(input))
	err := decoder.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("deserializing block from gob '%s': %s", input, err)
	}

	return FromWireBlock(&data)
}

func (*blockGobEncoder) EncodeHeader(header *entity.BlockHeader) ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(ToWireBlockHeader(header))
	if err != nil {
		return nil, fmt.Errorf("serializing header to gob: %s", err)
	}

	return result.Bytes(), nil
}

func (*blockGobEncoder) DecodeHeader(input []byte) (*entity.BlockHeader, error) {
	var data BlockHeader

	decoder := gob.NewDecoder(bytes.NewReader(input))
	err := decoder.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("deserializing header from gob '%s': %s", input, err)
	}

	return FromWireBlockHeader(&data)
}
