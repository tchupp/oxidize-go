package actions

import (
	"testing"
	"os"
	"time"
	"github.com/boltdb/bolt"
	"github.com/tclchiam/block_n_go/blockchain"
	"bytes"
)

func TestNewBlockAction_Execute(t *testing.T) {
	const testDbFileName = "test_blockchain.db"
	const testBlocksBucket = "test_blocks"

	t.Run("Test", func(t *testing.T) {
		newBlockchainAction := NewBlockchainAction{
			dbFileName:   testDbFileName,
			blocksBucket: testBlocksBucket,
		}

		_, err := newBlockchainAction.Execute()
		if err != nil {
			t.Fatalf("NewBlockchainAction failed: %s", err)
		}

		genesisBlock, err := getLatestBlockTest(t, testDbFileName, testBlocksBucket)
		if err != nil {
			t.Fatalf("error: %s", err)
		}

		const newBlockData = "Send Theo 3 BTC"

		newBlockAction := NewBlockAction{
			dbFileName: testDbFileName,
			bucketName: testBlocksBucket,
			data:       newBlockData,
		}

		_, err = newBlockAction.Execute()
		if err != nil {
			t.Fatalf("NewBlockAction failed: %s", err)
		}

		newBlock, err := getLatestBlockTest(t, testDbFileName, testBlocksBucket)
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
			t.Fatalf("New block has bad Index, expected [%s], but was [%s]", 1, newBlock.Index)
		}

	})

	err := os.Remove(testDbFileName)
	if err != nil {
		t.Fatalf("deleting test db file: %s", err)
	}
}

func getLatestBlockTest(t *testing.T, dbFileName string, bucketName string) (*blockchain.CommittedBlock, error) {
	db, err := bolt.Open(dbFileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		t.Fatalf("Error opening db: %s", err)
	}

	defer db.Close()

	return getLatestBlock(db, []byte(bucketName))
}
