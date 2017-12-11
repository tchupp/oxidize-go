package blockchain

import (
	"github.com/boltdb/bolt"
)

type Iterator struct {
	current    *CommittedBlock
	nodeName   string
	bucketName []byte
}

func (it *Iterator) Next() (*Iterator, error) {
	db, err := openDB(it.nodeName)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var block *CommittedBlock

	err = db.View(func(tx *bolt.Tx) error {
		bucket, err := bucket(tx, it.bucketName)
		if err != nil {
			return err
		}

		block, err = ReadBlock(bucket, it.current.PreviousHash)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Iterator{current: block, nodeName: it.nodeName, bucketName: it.bucketName}, nil
}

func (it *Iterator) Index() int {
	return it.current.Index
}

func (it *Iterator) PreviousHash() []byte {
	return it.current.PreviousHash
}

func (it *Iterator) Data() []byte {
	return it.current.Data
}

func (it *Iterator) Hash() []byte {
	return it.current.Hash
}

func (it *Iterator) Nonce() int {
	return it.current.Nonce
}

func (it *Iterator) Validate() bool {
	return it.current.Validate()
}
