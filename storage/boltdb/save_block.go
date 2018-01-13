package boltdb

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/tclchiam/block_n_go/blockchain/entity"
)

func (r *blockBoltRepository) SaveBlock(block *entity.Block) error {
	db, err := openDB(r.name)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := bucket(tx, blocksBucketName)
		if err != nil {
			return err
		}

		err = writeBlock(bucket, block, r.blockEncoder)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func writeBlock(bucket *bolt.Bucket, block *entity.Block, encoder entity.BlockEncoder) error {
	blockData, err := encoder.EncodeBlock(block)
	if err != nil {
		return err
	}

	err = bucket.Put(block.Hash().Slice(), blockData)
	if err != nil {
		return fmt.Errorf("writing block: %s", err)
	}

	err = bucket.Put(latestBlockHashKey, block.Hash().Slice())
	if err != nil {
		return fmt.Errorf("writing last hash: %s", err)
	}

	return nil
}
