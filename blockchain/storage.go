package blockchain

import (
	"fmt"
	"github.com/boltdb/bolt"
)

var (
	latestBlockHashKey = []byte("l")
)

func WriteBlock(bucket *bolt.Bucket, block *CommittedBlock) error {
	blockData, err := block.Serialize()
	if err != nil {
		return fmt.Errorf("serializing block: %s", err)
	}

	err = bucket.Put(block.Hash, blockData)
	if err != nil {
		return fmt.Errorf("writing block: %s", err)
	}

	err = bucket.Put(latestBlockHashKey, block.Hash)
	if err != nil {
		return fmt.Errorf("writing last hash: %s", err)
	}

	return nil
}

func ReadLatestHash(bucket *bolt.Bucket) ([]byte, error) {
	latestBlockHash := bucket.Get(latestBlockHashKey)
	if latestBlockHash == nil {
		return nil, fmt.Errorf("could not find latest block hash")
	}

	return latestBlockHash, nil
}

func ReadBlock(bucket *bolt.Bucket, blockHash []byte) (*CommittedBlock, error) {
	latestBlockData := bucket.Get(blockHash)
	if latestBlockData == nil || len(latestBlockData) == 0 {
		return nil, fmt.Errorf("latest block data is empty: '%s'", latestBlockData)
	}

	block, err := DeserializeBlock(latestBlockData)
	if err != nil {
		return nil, fmt.Errorf("deserializing block '%s': %s", latestBlockData, err)
	}

	return block, err
}
