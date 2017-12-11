package blockchain

import (
	"testing"
	"os"
	"time"
	"github.com/boltdb/bolt"
	"bytes"
	"github.com/google/go-cmp/cmp"
)

func TestNewBlock(t *testing.T) {
	const testDbFileName = "test_blockchain.db"
	testBlocksBucketName := []byte("test_blocks")

	db, err := bolt.Open(testDbFileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		t.Fatalf("Error opening db: %s", err)
	}

	defer db.Close()

	t.Run("Test", func(t *testing.T) {
		_, err := open(db, []byte(testBlocksBucketName))
		if err != nil {
			t.Fatalf("OpenBlockchain failed: %s", err)
		}

		genesisBlock, err := readLatestBlock(db, testBlocksBucketName)
		if err != nil {
			t.Fatalf("error: %s", err)
		}

		const newBlockData = "Send Theo 3 BTC"

		head, err := newBlock(db, testBlocksBucketName, newBlockData)
		if err != nil {
			t.Fatalf("NewBlockAction failed: %s", err)
		}

		newBlock, err := readLatestBlock(db, testBlocksBucketName)
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
		if !cmp.Equal(head, newBlock) {
			t.Fatalf("Resulting block does not equal the latest block: %s", cmp.Diff(head, newBlock))
		}
	})

	err = os.Remove(testDbFileName)
	if err != nil {
		t.Fatalf("deleting test db file: %s", err)
	}
}
