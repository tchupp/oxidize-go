package boltdb

import (
	"github.com/boltdb/bolt"
	"github.com/tclchiam/block_n_go/blockchain"
	"github.com/tclchiam/block_n_go/blockchain/chainhash"
	"github.com/tclchiam/block_n_go/blockchain/entity"
)

func (r *blockReader) Head() (head *entity.Block, err error) {
	err = r.db.View(func(tx *bolt.Tx) error {
		bucket, err := bucket(tx, blocksBucketName)
		if err != nil {
			return err
		}

		latestBlockHash := readLatestHash(bucket)
		if latestBlockHash == nil {
			return blockchain.HeadBlockNotFoundError
		}

		head, err = readBlock(bucket, latestBlockHash, r.blockEncoder)
		if err != nil {
			return err
		}

		return nil
	})

	return head, err
}

func (r *blockReader) Block(hash chainhash.Hash) (block *entity.Block, err error) {
	err = r.db.View(func(tx *bolt.Tx) error {
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
		return nil, blockchain.BlockDataEmptyError
	}

	b, err := encoder.DecodeBlock(latestBlockData)
	if err != nil {
		return nil, err
	}

	return b, err
}
