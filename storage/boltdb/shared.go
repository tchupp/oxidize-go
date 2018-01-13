package boltdb

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

func bucket(tx *bolt.Tx, bucketName []byte) (*bolt.Bucket, error) {
	bucket := tx.Bucket(bucketName)
	if bucket == nil {
		return nil, BucketNotFoundError
	}
	return bucket, nil
}

func openDB(name string) (*bolt.DB, error) {
	path := fmt.Sprintf(dbFile, name)
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("opening db: %s", err)
	}
	return db, err
}
