package boltdb

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/tclchiam/block_n_go/blockchain/block"
)

func (repo *BlockchainRepository) SaveBlock(block *block.Block) error {
	err := repo.db.Update(func(tx *bolt.Tx) error {
		bucket, err := bucket(tx, blocksBucketName)
		if err != nil {
			return err
		}

		err = writeBlock(bucket, block)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func writeBlock(bucket *bolt.Bucket, block *block.Block) error {
	blockData, err := SerializeBlock(block)
	if err != nil {
		return err
	}

	err = bucket.Put(block.Hash.Slice(), blockData)
	if err != nil {
		return fmt.Errorf("writing block: %s", err)
	}

	err = bucket.Put(latestBlockHashKey, block.Hash.Slice())
	if err != nil {
		return fmt.Errorf("writing last hash: %s", err)
	}

	return nil
}
