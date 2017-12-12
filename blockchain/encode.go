package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func (block *Block) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil {
		return nil, fmt.Errorf("serializing block: %s", err)
	}

	return result.Bytes(), nil
}

func DeserializeBlock(data []byte) (*Block, error) {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		return nil, fmt.Errorf("deserializing block '%s': %s", data, err)
	}

	return &block, nil
}
