package blockchain

import (
	"github.com/boltdb/bolt"
)

type Iterator struct {
	current    *Block
	nodeName   string
	bucketName []byte
}

func (it *Iterator) Next() (*Iterator, error) {
	db, err := openDB(it.nodeName)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var block *Block

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

	return &Iterator{
		current:    block,
		nodeName:   it.nodeName,
		bucketName: it.bucketName,
	}, nil
}

func (it *Iterator) HasNext() bool {
	return !it.current.IsGenesisBlock()
}

func (bc *Blockchain) Head() *Iterator {
	return &Iterator{
		current:    bc.head,
		nodeName:   bc.nodeName,
		bucketName: []byte(blockBucketName),
	}
}

func (bc *Blockchain) ForEachBlock(consume func(*Block)) (err error) {
	block := bc.Head()

	for {
		consume(block.current)

		if !block.HasNext() {
			break
		}

		block, err = block.Next()
		if err != nil {
			return err
		}
	}
	return nil
}
