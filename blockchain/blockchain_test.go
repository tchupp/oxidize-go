package blockchain

import (
	"fmt"
	"github.com/boltdb/bolt"
)

func readLatestBlock(db *bolt.DB, bucketName []byte) (latestBlock *CommittedBlock, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("no block with name %s exists", bucketName)
		}

		latestBlockHash := bucket.Get([]byte("l"))
		if latestBlockHash == nil {
			return fmt.Errorf("could not find latest block hash")
		}

		latestBlockData := bucket.Get(latestBlockHash)
		if latestBlockData == nil || len(latestBlockData) == 0 {
			return fmt.Errorf("latest block data is empty: '%s'", latestBlockData)
		}

		latestBlock, err = DeserializeBlock(latestBlockData)
		if err != nil {
			return err
		}

		return nil
	})

	return latestBlock, err
}
