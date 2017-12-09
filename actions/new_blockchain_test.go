package actions

import (
	"testing"
	"github.com/boltdb/bolt"
	"github.com/tclchiam/block_n_go/blockchain"
	"time"
	"os"
)

func TestNewBlockchainAction_Execute(t *testing.T) {
	const testBlocksBucket = "test_blocks"
	const testDbFileName = "test_blockchain.db"

	t.Run("Test", func(t *testing.T) {
		action := NewBlockchainAction{
			dbFileName:   testDbFileName,
			blocksBucket: testBlocksBucket,
		}

		_, err := action.Execute()
		if err != nil {
			t.Fatalf("NewBlockchainAction failed: %s", err)
		}

		db, err := bolt.Open(testDbFileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			t.Fatalf("Error opening db: %s", err)
		}

		defer db.Close()

		err = db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(testBlocksBucket))
			if bucket == nil {
				t.Fatalf("no block with name '%s' exists", testBlocksBucket)
			}

			lastHash := bucket.Get([]byte("l"))
			if lastHash == nil {
				t.Fatalf("could not find last hash")
			}

			blockData := bucket.Get(lastHash)
			if blockData == nil || len(blockData) == 0 {
				t.Fatalf("block data is empty: '%s'", blockData)
			}

			genesisBlock, err := blockchain.DeserializeBlock(blockData)
			if err != nil {
				t.Fatalf("deserializing block '%s': %s", blockData, err)
			}
			if genesisBlock == nil {
				t.Fatalf("Genesis block is nil")
			}
			if len(genesisBlock.PreviousHash) != 0 {
				t.Fatalf("Genesis block has bad PreviousHash, expected [%s], but was [%s]", []byte{}, genesisBlock.PreviousHash)
			}

			return nil
		})
		if err != nil {
			t.Fatalf("error: %s", err)
		}
	})

	err := os.Remove(testDbFileName)
	if err != nil {
		t.Fatalf("deleting test db file: %s", err)
	}
}
