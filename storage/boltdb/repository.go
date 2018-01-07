package boltdb

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain_%s.db"

var (
	latestBlockHashKey = []byte("l")
	blocksBucketName   = []byte("blocks")
)

type BlockchainRepository struct {
	// Filename to the BoltDB database
	Path string

	db *bolt.DB
}

func NewRepository(name string) (*BlockchainRepository, error) {
	path := fmt.Sprintf(dbFile, name)

	db, err := openDB(path)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		return createBucket(tx, blocksBucketName)
	})
	if err != nil {
		db.Close()
		return nil, err
	}

	return &BlockchainRepository{Path: path, db: db}, nil

}

func openDB(dbFile string) (*bolt.DB, error) {
	db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("opening db: %s", err)
	}
	return db, err
}

func (repo *BlockchainRepository) Close() error {
	return repo.db.Close()
}

func createBucket(tx *bolt.Tx, bucketName []byte) error {
	if _, err := tx.CreateBucketIfNotExists(bucketName); err != nil {
		return fmt.Errorf("creating block bucket: %s", err)
	}
	return nil
}
