package actions

import (
	"github.com/boltdb/bolt"
	"github.com/tclchiam/block_n_go/blockchain"
	"fmt"
)

type NewBlockAction struct {
	data       string
	bucketName []byte
}

func (action *NewBlockAction) Execute(db *bolt.DB) (*blockchain.Blockchain, error) {
	latestBlock, err := getLatestBlock(db, action.bucketName)
	if err != nil {
		return nil, fmt.Errorf("reading last hash: %s", err)
	}

	newBlock := blockchain.NewBlock(action.data, latestBlock.Hash, latestBlock.Index+1)

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(action.bucketName)
		if bucket == nil {
			return fmt.Errorf("no block with name '%s' exists", action.bucketName)
		}

		err = blockchain.WriteBlock(bucket, newBlock)
		if err != nil {
			return err
		}

		return nil
	})

	return blockchain.New(newBlock.Hash), nil
}

func getLatestBlock(db *bolt.DB, bucketNameBytes []byte) (*blockchain.CommittedBlock, error) {
	var latestBlock *blockchain.CommittedBlock
	var err error

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketNameBytes)
		if bucket == nil {
			return fmt.Errorf("no block with name %s exists", bucketNameBytes)
		}

		latestBlockHash, err := blockchain.ReadLatestHash(bucket)
		if err != nil {
			return err
		}

		latestBlock, err = blockchain.ReadBlock(bucket, latestBlockHash)
		if err != nil {
			return err
		}

		return nil
	})

	return latestBlock, err
}
