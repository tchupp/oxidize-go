package actions

import (
	"testing"
	"github.com/boltdb/bolt"
	"github.com/tclchiam/block_n_go/blockchain"
	"time"
	"os"
)

func TestNewBlockchainAction_Execute(t *testing.T) {
	const testDbFileName = "test_blockchain.db"
	const testBlocksBucket = "test_blocks"

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

			latestBlockHash := bucket.Get([]byte("l"))
			if latestBlockHash == nil {
				t.Fatalf("could not find latest block hash")
			}

			latestBlockData := bucket.Get(latestBlockHash)
			if latestBlockData == nil || len(latestBlockData) == 0 {
				t.Fatalf("latest block data is empty: '%s'", latestBlockData)
			}

			genesisBlock, err := blockchain.DeserializeBlock(latestBlockData)
			if err != nil {
				t.Fatalf("deserializing block '%s': %s", latestBlockData, err)
			}

			if len(genesisBlock.PreviousHash) != 0 {
				t.Fatalf("Genesis block has bad PreviousHash, expected [%s], but was [%s]", []byte{}, genesisBlock.PreviousHash)
			}
			if genesisBlock.Index != 0 {
				t.Fatalf("Genesis block has bad Index, expected [%s], but was [%s]", 0, genesisBlock.Index)
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
