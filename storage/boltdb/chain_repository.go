package boltdb

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"

	"github.com/tclchiam/block_n_go/blockchain/entity"
	"os"
)

const dbFile = "blockchain_%s.db"

const (
	blocksBucketName  = "blocks"
	headersBucketName = "headers"

	indexKeyLength = 8
)

type chainBoltRepository struct {
	name         string
	blockEncoder entity.BlockEncoder
	db           *bolt.DB
}

func NewChainRepository(name string, blockEncoder entity.BlockEncoder) (entity.ChainRepository, error) {
	db, err := openDB(name)
	if err != nil {
		return nil, err
	}

	if err = createBucket(db, blocksBucketName); err != nil {
		return nil, err
	}
	if err = createBucket(db, headersBucketName); err != nil {
		return nil, err
	}

	return &chainBoltRepository{
		name:         name,
		blockEncoder: blockEncoder,
		db:           db,
	}, nil
}

func NewBlockRepository(name string, blockEncoder entity.BlockEncoder) (entity.BlockRepository, error) {
	return NewChainRepository(name, blockEncoder)
}

func NewHeaderRepository(name string, blockEncoder entity.BlockEncoder) (entity.BlockRepository, error) {
	return NewChainRepository(name, blockEncoder)
}

func DeleteBlockchain(name string) error {
	path := fmt.Sprintf(dbFile, name)
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("deleting blockchain file: %s", err)
	}
	return nil
}

func (r *chainBoltRepository) BestBlock() (block *entity.Block, err error) {
	err = r.db.View(func(tx *bolt.Tx) error {
		bestHash, err := bestBlockHash(tx)
		if err != nil {
			return err
		}
		if bestHash == nil {
			return nil
		}

		block, err = blockByHash(tx, r.blockEncoder, bestHash)
		return err
	})
	return block, err
}

func (r *chainBoltRepository) BlockByHash(hash *entity.Hash) (block *entity.Block, err error) {
	err = r.db.View(func(tx *bolt.Tx) error {
		block, err = blockByHash(tx, r.blockEncoder, hash)
		return err
	})
	return block, err
}

func (r *chainBoltRepository) BlockByIndex(index uint64) (block *entity.Block, err error) {
	err = r.db.View(func(tx *bolt.Tx) error {
		hash, err := hashByIndex(tx, index)
		if err != nil {
			return err
		}
		if hash == nil {
			return nil
		}

		block, err = blockByHash(tx, r.blockEncoder, hash)
		return err
	})
	return block, err
}

func (r *chainBoltRepository) SaveBlock(block *entity.Block) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		return saveBlock(tx, block, r.blockEncoder)
	})
}

func (r *chainBoltRepository) BestHeader() (header *entity.BlockHeader, err error) {
	err = r.db.View(func(tx *bolt.Tx) error {
		bestHash, err := bestHeaderHash(tx)
		if err != nil {
			return err
		}
		if bestHash == nil {
			return nil
		}

		header, err = headerByHash(tx, r.blockEncoder, bestHash)
		return err
	})
	return header, err
}

func (r *chainBoltRepository) HeaderByHash(hash *entity.Hash) (header *entity.BlockHeader, err error) {
	err = r.db.View(func(tx *bolt.Tx) error {
		header, err = headerByHash(tx, r.blockEncoder, hash)
		return err
	})
	return header, err
}

func (r *chainBoltRepository) HeaderByIndex(index uint64) (header *entity.BlockHeader, err error) {
	err = r.db.View(func(tx *bolt.Tx) error {
		hash, err := hashByIndex(tx, index)
		if err != nil {
			return err
		}
		if hash == nil {
			return nil
		}

		header, err = headerByHash(tx, r.blockEncoder, hash)
		return err
	})
	return header, err
}

func (r *chainBoltRepository) SaveHeader(header *entity.BlockHeader) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		return saveHeader(tx, r.blockEncoder, header)
	})
}

func (r *chainBoltRepository) Close() error {
	return r.db.Close()
}

func openDB(name string) (*bolt.DB, error) {
	path := fmt.Sprintf(dbFile, name)
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("opening db: %s", err)
	}
	return db, err
}

func hashByIndex(tx *bolt.Tx, index uint64) (*entity.Hash, error) {
	bucket, err := bucket(tx, blocksBucketName)
	if err != nil {
		return nil, err
	}

	if rawHash := bucket.Get(toByte(index)); rawHash != nil {
		return entity.NewHash(rawHash)
	}
	return nil, nil
}
