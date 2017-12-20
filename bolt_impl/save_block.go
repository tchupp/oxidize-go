package bolt_impl

import (
	"github.com/tclchiam/block_n_go/blockchain"
	"github.com/boltdb/bolt"
	"fmt"
)

func (repo *BlockchainRepository) SaveBlock(block *blockchain.Block) error {
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

func writeBlock(bucket *bolt.Bucket, block *blockchain.Block) error {
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
