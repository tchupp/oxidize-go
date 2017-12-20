package bolt_impl

import (
	"github.com/boltdb/bolt"
	"github.com/tclchiam/block_n_go/blockchain"
)

func bucket(tx *bolt.Tx, bucketName []byte) (*bolt.Bucket, error) {
	bucket := tx.Bucket(bucketName)
	if bucket == nil {
		return nil, blockchain.BucketNotFoundError
	}
	return bucket, nil
}