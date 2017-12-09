package actions

import (
	"testing"
	"github.com/boltdb/bolt"
	"github.com/tclchiam/block_n_go/blockchain"
	"time"
)

func TestNewBlockchainAction_Execute(t *testing.T) {
	const testBlocksBucket = "test_blocks"
	const testDbFileName = "test_blockchain.db"

	action := NewBlockchainAction{
		dbFileName:   testDbFileName,
		blocksBucket: testBlocksBucket,
	}

	result, err := action.Execute()

	if !result {
		t.Errorf("NewBlockchainAction failed")
	}

	db, err := bolt.Open(testDbFileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		t.Errorf("Error opening db: %s", err)
	}

	defer db.Close()

	var genesisBlock blockchain.CommittedBlock
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(testBlocksBucket))

		if b == nil {
			t.Errorf("Block %s is nil, expected to exist", testBlocksBucket)
		}

		lastHash := b.Get([]byte("l"))
		blockData := b.Get(lastHash)
		genesisBlock = *blockchain.DeserializeBlock(blockData)

		return nil
	})

	if len(genesisBlock.PreviousHash) != 0 {
		t.Errorf("Genesis block has bad PreviousHash, expected [%s], but was [%s]", []byte{}, genesisBlock.PreviousHash)
	}
}
