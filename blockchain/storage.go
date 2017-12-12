package blockchain

import (
	"fmt"
	"github.com/boltdb/bolt"
)

var (
	latestBlockHashKey = []byte("l")
)

func WriteBlock(bucket *bolt.Bucket, block *Block) error {
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

func ReadBlock(bucket *bolt.Bucket, blockHash []byte) (*Block, error) {
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

func ReadHeadBlock(db *bolt.DB, bucketName []byte) (headBlock *Block, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		bucket, err := bucket(tx, bucketName)
		if err != nil {
			return err
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

func bucket(tx *bolt.Tx, bucketName []byte) (*bolt.Bucket, error) {
	bucket := tx.Bucket(bucketName)
	if bucket == nil {
		return nil, fmt.Errorf("no bucket with name '%s' exists", bucketName)
	}
	return bucket, nil
}
