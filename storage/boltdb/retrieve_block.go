package boltdb

import (
	"github.com/boltdb/bolt"
	"github.com/tclchiam/block_n_go/blockchain/entity"
)

func (r *blockBoltRepository) Head() (head *entity.Block, err error) {
	db, err := openDB(r.name)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		bucket, err := bucket(tx, blocksBucketName)
		if err != nil {
			return err
		}

		latestBlockHash := readLatestHash(bucket)
		if latestBlockHash == nil {
			return nil
		}

		head, err = readBlock(bucket, latestBlockHash, r.blockEncoder)
		if err != nil {
			return err
		}

		return nil
	})

	return head, err
}

func (r *blockBoltRepository) Block(hash *entity.Hash) (block *entity.Block, err error) {
	db, err := openDB(r.name)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		bucket, err := bucket(tx, blocksBucketName)
		if err != nil {
			return err
		}

		block, err = readBlock(bucket, hash.Slice(), r.blockEncoder)
		if err != nil {
			return err
		}

		return nil
	})

	return block, err
}

func readLatestHash(bucket *bolt.Bucket) []byte {
	return bucket.Get(latestBlockHashKey)
}

func readBlock(bucket *bolt.Bucket, hash []byte, encoder entity.BlockEncoder) (*entity.Block, error) {
	latestBlockData := bucket.Get(hash)
	if latestBlockData == nil || len(latestBlockData) == 0 {
		return nil, BlockDataEmptyError
	}

	b, err := encoder.DecodeBlock(latestBlockData)
	if err != nil {
		return nil, err
	}

	return b, err
}
