package boltdb

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/tclchiam/block_n_go/blockchain/entity"
)

const dbFile = "blockchain_%s.db"

var (
	latestBlockHashKey = []byte("l")
	blocksBucketName   = []byte("blocks")
)

type blockBoltRepository struct {
	name         string
	blockEncoder entity.BlockEncoder
}

func NewBlockRepository(name string, blockEncoder entity.BlockEncoder) (entity.BlockRepository, error) {
	db, err := openDB(name)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if err = createBucket(db, blocksBucketName); err != nil {
		return nil, err
	}

	return &blockBoltRepository{
		name:         name,
		blockEncoder: blockEncoder,
	}, nil
}

func createBucket(db *bolt.DB, bucketName []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(bucketName); err != nil {
			return fmt.Errorf("creating block bucket: %s", err)
		}
		return nil
	})
}
