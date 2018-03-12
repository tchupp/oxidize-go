package boltdb

import (
	"fmt"
	"time"

	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
)

const (
	chainDbFile = "blockchain_%s.db"

	blocksBucketName  = "blocks"
	headersBucketName = "headers"

	indexKeyLength = 8
)

type chainBoltRepository struct {
	name string
	db   *bolt.DB
}

func NewChainRepository(path, name string) (entity.ChainRepository, error) {
	db, err := openDB(path, fmt.Sprintf(chainDbFile, name))
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
		db: db,
	}, nil
}

func DeleteBlockchain(name string) error {
	return DeleteBoltFile("./", fmt.Sprintf(chainDbFile, name))
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

		block, err = blockByHash(tx, bestHash)
		return err
	})
	return block, err
}

func (r *chainBoltRepository) BlockByHash(hash *entity.Hash) (block *entity.Block, err error) {
	err = r.db.View(func(tx *bolt.Tx) error {
		block, err = blockByHash(tx, hash)
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

		block, err = blockByHash(tx, hash)
		return err
	})
	return block, err
}

func (r *chainBoltRepository) SaveBlock(block *entity.Block) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		return saveBlock(tx, block)
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

		header, err = headerByHash(tx, bestHash)
		return err
	})
	return header, err
}

func (r *chainBoltRepository) HeaderByHash(hash *entity.Hash) (header *entity.BlockHeader, err error) {
	err = r.db.View(func(tx *bolt.Tx) error {
		header, err = headerByHash(tx, hash)
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

		header, err = headerByHash(tx, hash)
		return err
	})
	return header, err
}

func (r *chainBoltRepository) SaveHeader(header *entity.BlockHeader) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		return saveHeader(tx, header)
	})
}

func (r *chainBoltRepository) Close() error {
	return r.db.Close()
}

func openDB(path, name string) (*bolt.DB, error) {
	path = filepath.Join(path, name)
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
