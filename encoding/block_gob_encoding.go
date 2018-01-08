package encoding

import (
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"bytes"
	"encoding/gob"
	"fmt"
)

type blockGobEncoder struct{}

var blockGobEncoderInstance blockGobEncoder

func NewBlockGobEncoder() entity.BlockEncoder {
	return &blockGobEncoderInstance
}

func (*blockGobEncoder) EncodeBlock(block *entity.Block) ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(toBlockData(block))
	if err != nil {
		return nil, fmt.Errorf("serializing block to gob: %s", err)
	}

	return result.Bytes(), nil
}

func (*blockGobEncoder) DecodeBlock(input []byte) (*entity.Block, error) {
	var data blockData

	decoder := gob.NewDecoder(bytes.NewReader(input))
	err := decoder.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("deserializing block from gob '%s': %s", input, err)
	}

	return fromBlockData(&data)
}
