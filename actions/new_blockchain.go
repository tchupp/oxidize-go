package actions

import (
	"github.com/boltdb/bolt"
	"github.com/tclchiam/block_n_go/blockchain"
	"fmt"
)

type NewBlockchainAction struct {
	bucketName []byte
}

func (action *NewBlockchainAction) Execute(db *bolt.DB) (bool, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		blocksBucketExists := tx.Bucket(action.bucketName) != nil

		if blocksBucketExists == false {
			genesisBlock := blockchain.NewGenesisBlock()

			bucket, err := tx.CreateBucket(action.bucketName)
			if err != nil {
				return fmt.Errorf("creating block bucket: %s", err)
			}

			err = blockchain.WriteBlock(bucket, genesisBlock)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return false, err
	}

	return true, nil
}
