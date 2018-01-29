package boltdb

import (
	"fmt"
	"os"
	"time"

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

func DeleteBlockchain(name string) error {
	path := fmt.Sprintf(dbFile, name)
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("deleting blockchain file: %s", err)
	}
	return nil
}

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

func (r *blockBoltRepository) SaveBlock(block *entity.Block) error {
	db, err := openDB(r.name)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := bucket(tx, blocksBucketName)
		if err != nil {
			return err
		}

		err = writeBlock(bucket, block, r.blockEncoder)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func writeBlock(bucket *bolt.Bucket, block *entity.Block, encoder entity.BlockEncoder) error {
	blockData, err := encoder.EncodeBlock(block)
	if err != nil {
		return err
	}

	err = bucket.Put(block.Hash().Slice(), blockData)
	if err != nil {
		return fmt.Errorf("writing block: %s", err)
	}

	err = bucket.Put(latestBlockHashKey, block.Hash().Slice())
	if err != nil {
		return fmt.Errorf("writing last hash: %s", err)
	}

	return nil
}
