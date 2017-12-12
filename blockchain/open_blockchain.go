package blockchain

import (
	"github.com/boltdb/bolt"
	"fmt"
)

func open(db *bolt.DB, bucketName []byte, address string) (headBlock *Block, err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)

		if bucket == nil {
			genesisBlock := NewGenesisBlock(address)

			bucket, err := tx.CreateBucket(bucketName)
			if err != nil {
				return fmt.Errorf("creating block bucket: %s", err)
			}

			err = WriteBlock(bucket, genesisBlock)
			if err != nil {
				return err
			}

			headBlock = genesisBlock
		} else {
			headBlock, err = ReadHeadBlock(db, bucketName)
			if err != nil {
				return err
			}
		}

		return nil
	})
	return headBlock, err
}
