package bolt_impl

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/tclchiam/block_n_go/blockchain"
)

func SerializeBlock(block *blockchain.Block) ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil {
		return nil, fmt.Errorf("serializing block: %s", err)
	}

	return result.Bytes(), nil
}

func DeserializeBlock(data []byte) (*blockchain.Block, error) {
	var block blockchain.Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		return nil, fmt.Errorf("deserializing block '%s': %s", data, err)
	}

	return &block, nil
}
