package actions

import (
	"testing"
	"os"
	"time"
	"github.com/boltdb/bolt"
	"github.com/tclchiam/block_n_go/blockchain"
	"bytes"
	"fmt"
)

func TestNewBlockAction_Execute(t *testing.T) {
	const testDbFileName = "test_blockchain.db"
	const testBlocksBucketName = "test_blocks"

	db, err := bolt.Open(testDbFileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		t.Fatalf("Error opening db: %s", err)
	}

	defer db.Close()

	t.Run("Test", func(t *testing.T) {
		newBlockchainAction := NewBlockchainAction{
			bucketName: testBlocksBucketName,
		}

		_, err := newBlockchainAction.Execute(db)
		if err != nil {
			t.Fatalf("NewBlockchainAction failed: %s", err)
		}

		genesisBlock, err := testGetLatestBlock(db, []byte(testBlocksBucketName))
		if err != nil {
			t.Fatalf("error: %s", err)
		}

		const newBlockData = "Send Theo 3 BTC"

		newBlockAction := NewBlockAction{
			bucketName: testBlocksBucketName,
			data:       newBlockData,
		}

		_, err = newBlockAction.Execute(db)
		if err != nil {
			t.Fatalf("NewBlockAction failed: %s", err)
		}

		newBlock, err := testGetLatestBlock(db, []byte(testBlocksBucketName))
		if err != nil {
			t.Fatalf("error: %s", err)
		}

		if bytes.Compare(newBlock.PreviousHash, genesisBlock.Hash) != 0 {
			t.Fatalf("New block has bad PreviousHash, expected [%s], but was [%s]", genesisBlock.Hash, newBlock.PreviousHash)
		}
		if newBlock.Index != 1 {
			t.Fatalf("New block has bad Index, expected [%s], but was [%s]", 1, newBlock.Index)
		}
		if string(newBlock.Data) != newBlockData {
			t.Fatalf("New block has bad Index, expected [%s], but was [%s]", newBlockData, newBlock.Data)
		}

	})

	err = os.Remove(testDbFileName)
	if err != nil {
		t.Fatalf("deleting test db file: %s", err)
	}
}

func testGetLatestBlock(db *bolt.DB, bucketNameBytes []byte) (*blockchain.CommittedBlock, error) {
	var latestBlock *blockchain.CommittedBlock
	var err error

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketNameBytes)
		if bucket == nil {
			return fmt.Errorf("no block with name %s exists", bucketNameBytes)
		}

		latestBlockHash := bucket.Get([]byte("l"))
		if latestBlockHash == nil {
			return fmt.Errorf("could not find latest block hash")
		}

		latestBlockData := bucket.Get(latestBlockHash)
		if latestBlockData == nil || len(latestBlockData) == 0 {
			return fmt.Errorf("latest block data is empty: '%s'", latestBlockData)
		}

		latestBlock, err = blockchain.DeserializeBlock(latestBlockData)
		if err != nil {
			return fmt.Errorf("deserializing block '%s': %s", latestBlockData, err)
		}

		return nil
	})

	return latestBlock, err
}
