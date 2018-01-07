package boltdb

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/tclchiam/block_n_go/blockchain/block"
)

func SerializeBlock(block *block.Block) ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil {
		return nil, fmt.Errorf("serializing block: %s", err)
	}

	return result.Bytes(), nil
}

func DeserializeBlock(data []byte) (*block.Block, error) {
	var b block.Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&b)
	if err != nil {
		return nil, fmt.Errorf("deserializing block '%s': %s", data, err)
	}

	return &b, nil
}
