package boltdb

import (
	"github.com/boltdb/bolt"
)

func bucket(tx *bolt.Tx, bucketName []byte) (*bolt.Bucket, error) {
	bucket := tx.Bucket(bucketName)
	if bucket == nil {
		return nil, BucketNotFoundError
	}
	return bucket, nil
}
