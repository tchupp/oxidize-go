package actions

import (
	"github.com/boltdb/bolt"
	"github.com/tclchiam/block_n_go/blockchain"
	"fmt"
)

type NewBlockAction struct {
	data       string
	bucketName string
}

func (action *NewBlockAction) Execute(db *bolt.DB) (bool, error) {
	bucketNameBytes := []byte(action.bucketName)

	latestBlock, err := getLatestBlock(db, bucketNameBytes)
	if err != nil {
		return false, fmt.Errorf("reading last hash: %s", err)
	}

	newBlock := blockchain.NewBlock(action.data, latestBlock.Hash, latestBlock.Index+1)

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketNameBytes)
		if bucket == nil {
			return fmt.Errorf("no block with name '%s' exists", bucketNameBytes)
		}

		blockData, err := newBlock.Serialize()
		if err != nil {
			return fmt.Errorf("serializing block: %s", err)
		}

		err = bucket.Put(newBlock.Hash, blockData)
		if err != nil {
			return fmt.Errorf("writing block: %s", err)
		}

		err = bucket.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			return fmt.Errorf("writing last hash: %s", err)
		}

		return nil
	})

	return true, nil
}

func getLatestBlock(db *bolt.DB, bucketNameBytes []byte) (*blockchain.CommittedBlock, error) {
	var latestBlock *blockchain.CommittedBlock
	var err error

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketNameBytes)
		if bucket == nil {
			return fmt.Errorf("no block with name %s exists", bucketNameBytes)
		}

		latestBlockHash := bucket.Get([]byte("l"))
		if latestBlockHash == nil {
			return fmt.Errorf("could not find latest block hash")
		}

		latestBlockData := bucket.Get(latestBlockHash)
		if latestBlockData == nil || len(latestBlockData) == 0 {
			return fmt.Errorf("latest block data is empty: '%s'", latestBlockData)
		}

		latestBlock, err = blockchain.DeserializeBlock(latestBlockData)
		if err != nil {
			return fmt.Errorf("deserializing block '%s': %s", latestBlockData, err)
		}

		return nil
	})

	return latestBlock, err
}
