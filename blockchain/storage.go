package blockchain

import (
	"fmt"
	"github.com/boltdb/bolt"
	"bytes"
	"encoding/gob"
)

var (
	latestBlockHashKey = []byte("l")
)

func (block *CommittedBlock) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil {
		return nil, fmt.Errorf("serializing block: %s", err)
	}

	return result.Bytes(), nil
}

func DeserializeBlock(data []byte) (*CommittedBlock, error) {
	var block CommittedBlock

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		return nil, fmt.Errorf("deserializing block '%s': %s", data, err)
	}

	return &block, nil
}


func WriteBlock(bucket *bolt.Bucket, block *CommittedBlock) error {
	blockData, err := block.Serialize()
	if err != nil {
		return err
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

func ReadBlock(bucket *bolt.Bucket, blockHash []byte) (*CommittedBlock, error) {
	latestBlockData := bucket.Get(blockHash)
	if latestBlockData == nil || len(latestBlockData) == 0 {
		return nil, fmt.Errorf("block data is empty: '%s'", latestBlockData)
	}

	block, err := DeserializeBlock(latestBlockData)
	if err != nil {
		return nil, err
	}

	return block, err
}

func ReadLatestHash(bucket *bolt.Bucket) ([]byte, error) {
	latestBlockHash := bucket.Get(latestBlockHashKey)
	if latestBlockHash == nil {
		return nil, fmt.Errorf("could not find latest block hash")
	}

	return latestBlockHash, nil
}

func ReadHeadBlock(db *bolt.DB, bucketNameBytes []byte) (headBlock *CommittedBlock, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketNameBytes)
		if bucket == nil {
			return fmt.Errorf("no block with name %s exists", bucketNameBytes)
		}

		latestBlockHash, err := ReadLatestHash(bucket)
		if err != nil {
			return err
		}

		headBlock, err = ReadBlock(bucket, latestBlockHash)
		if err != nil {
			return err
		}

		return nil
	})

	return headBlock, err
}
