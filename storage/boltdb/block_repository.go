package boltdb

import (
	"bytes"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/wire"
)

func bestBlockHash(tx *bolt.Tx) (hash *entity.Hash, err error) {
	bucket, err := bucket(tx, blocksBucketName)
	if err != nil {
		return nil, err
	}

	index, err := bestIndex(bucket)
	if err != nil {
		return nil, err
	}
	if index == nil {
		return nil, err
	}

	rawHash := bucket.Get(index)
	return entity.NewHash(rawHash)
}

func blockByHash(tx *bolt.Tx, hash *entity.Hash) (*entity.Block, error) {
	bucket, err := bucket(tx, blocksBucketName)
	if err != nil {
		return nil, err
	}

	block, err := readBlock(bucket, hash)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func readBlock(bucket *bolt.Bucket, hash *entity.Hash) (*entity.Block, error) {
	latestBlockData := bucket.Get(hash.Slice())
	if latestBlockData == nil || len(latestBlockData) == 0 {
		return nil, nil
	}

	b, err := wire.DecodeBlock(latestBlockData)
	if err != nil {
		return nil, err
	}

	return b, err
}

func saveBlock(tx *bolt.Tx, block *entity.Block) error {
	bucket, err := bucket(tx, blocksBucketName)
	if err != nil {
		return err
	}

	blockData, err := wire.EncodeBlock(block)
	if err != nil {
		return fmt.Errorf("serializing block: %s", err)
	}

	err = bucket.Put(block.Hash().Slice(), blockData)
	if err != nil {
		return fmt.Errorf("writing block: %s", err)
	}

	err = bucket.Put(toByte(block.Index()), block.Hash().Slice())
	if err != nil {
		return fmt.Errorf("writing block hash: %s", err)
	}

	return saveHeader(tx, block.Header())
}

func bestIndex(bucket *bolt.Bucket) (index []byte, err error) {
	err = bucket.ForEach(func(k, v []byte) error {
		if len(k) == indexKeyLength && bytes.Compare(index, k) == -1 {
			index = k
		}
		return nil
	})
	return index, err
}
