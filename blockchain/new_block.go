package blockchain

import (
	"github.com/boltdb/bolt"
	"fmt"
	"github.com/tclchiam/block_n_go/tx"
)

func newBlock(db *bolt.DB, bucketName []byte, transactions []*tx.Transaction) (newBlock *Block, err error) {
	latestBlock, err := ReadHeadBlock(db, bucketName)
	if err != nil {
		return nil, fmt.Errorf("reading head block: %s", err)
	}

	newBlock = NewBlock(transactions, latestBlock.Hash, latestBlock.Index+1)

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := bucket(tx, bucketName)
		if err != nil {
			return err
		}

		err = WriteBlock(bucket, newBlock)
		if err != nil {
			return err
		}

		return nil
	})

	return newBlock, err
}
