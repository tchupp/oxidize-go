package boltdb

import (
	"fmt"

	"github.com/boltdb/bolt"
)

func createBucket(db *bolt.DB, bucketName string) error {
	return db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(bucketName)); err != nil {
			return fmt.Errorf("creating '%s' bucket: %s", bucketName, err)
		}
		return nil
	})
}

func bucket(tx *bolt.Tx, bucketName string) (*bolt.Bucket, error) {
	bucket := tx.Bucket([]byte(bucketName))
	if bucket == nil {
		return nil, BucketNotFoundError
	}
	return bucket, nil
}
