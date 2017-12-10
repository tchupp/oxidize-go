package actions

import (
	"testing"
	"github.com/boltdb/bolt"
	"time"
	"os"
	"bytes"
)

func TestOpenBlockchainAction_Execute(t *testing.T) {
	const testDbFileName = "test_blockchain.db"
	var testBlocksBucketName = []byte("test_blocks")

	db, err := bolt.Open(testDbFileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		t.Fatalf("Error opening db: %s", err)
	}

	defer db.Close()

	t.Run("Test", func(t *testing.T) {
		action := OpenBlockchainAction{
			bucketName: []byte(testBlocksBucketName),
		}

		bc, err := action.Execute(db)
		if err != nil {
			t.Fatalf("OpenBlockchainAction failed: %s", err)
		}

		genesisBlock, err := testGetLatestBlock(db, testBlocksBucketName)
		if err != nil {
			t.Fatalf("error: %s", err)
		}

		if len(genesisBlock.PreviousHash) != 0 {
			t.Fatalf("Genesis block has bad PreviousHash, expected '%x', but was '%x'", []byte{}, genesisBlock.PreviousHash)
		}
		if genesisBlock.Index != 0 {
			t.Fatalf("Genesis block has bad Index, expected '%s', but was '%s'", 0, genesisBlock.Index)
		}
		if bytes.Compare(bc.LatestHash(), genesisBlock.Hash) != 0 {
			t.Fatalf("Resulting blockchain's latest hash does not match block's hash: expected '%x', was '%x'", genesisBlock.Hash, bc.LatestHash())
		}
	})

	err = os.Remove(testDbFileName)
	if err != nil {
		t.Fatalf("deleting test db file: %s", err)
	}
}
