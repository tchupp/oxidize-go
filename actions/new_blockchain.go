package actions

import (
	"github.com/boltdb/bolt"
	"github.com/tclchiam/block_n_go/blockchain"
)

type NewBlockchainAction struct {
	dbFileName   string
	blocksBucket string
}

func (action *NewBlockchainAction) Execute() (bool, error) {
	db, err := bolt.Open(action.dbFileName, 0600, nil)
	if err != nil {
		return false, err
	}

	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(action.blocksBucket))

		if b == nil {
			genesis := blockchain.NewGenesisBlock()

			b, err := tx.CreateBucket([]byte(action.blocksBucket))
			if err != nil {
				return err
			}

			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				return err
			}

			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return false, err
	}

	return true, nil
}
