package actions

import (
	"github.com/boltdb/bolt"
	"github.com/tclchiam/block_n_go/blockchain"
	"fmt"
)

type OpenBlockchainAction struct {
	bucketName []byte
}

func (action *OpenBlockchainAction) Execute(db *bolt.DB) (*blockchain.Blockchain, error) {
	var latestHash []byte
	var err error

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(action.bucketName)

		if bucket == nil {
			genesisBlock := blockchain.NewGenesisBlock()

			bucket, err := tx.CreateBucket(action.bucketName)
			if err != nil {
				return fmt.Errorf("creating block bucket: %s", err)
			}

			err = blockchain.WriteBlock(bucket, genesisBlock)
			if err != nil {
				return err
			}

			latestHash = genesisBlock.Hash
		} else {
			latestHash, err = blockchain.ReadLatestHash(bucket)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return blockchain.New(latestHash), nil
}
