package boltdb

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/tclchiam/block_n_go/storage"
)

const dbFile = "blockchain_%s.db"

var (
	latestBlockHashKey = []byte("l")
	blocksBucketName   = []byte("blocks")
)

type blockReader struct {
	// Filename to the BoltDB database
	name string

	db *bolt.DB
}

func NewReader(name string) (storage.BlockReader, error) {
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

	return &blockReader{name: name, db: db}, nil

}

func openDB(dbFile string) (*bolt.DB, error) {
	db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("opening db: %s", err)
	}
	return db, err
}

func (r *blockReader) Close() error {
	return r.db.Close()
}

func createBucket(tx *bolt.Tx, bucketName []byte) error {
	if _, err := tx.CreateBucketIfNotExists(bucketName); err != nil {
		return fmt.Errorf("creating block bucket: %s", err)
	}
	return nil
}
