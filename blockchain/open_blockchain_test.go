package blockchain

import (
	"testing"
	"time"
	"os"
	"github.com/boltdb/bolt"
	"github.com/google/go-cmp/cmp"
)

func TestOpen(t *testing.T) {
	const address = "491823"
	const testDbFileName = "test_blockchain.db"
	var testBlocksBucketName = []byte("test_blocks")

	db, err := bolt.Open(testDbFileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		t.Fatalf("Error opening db: %s", err)
	}

	defer db.Close()

	t.Run("Test", func(t *testing.T) {
		head, err := open(db, []byte(testBlocksBucketName), address)

		if err != nil {
			t.Fatalf("OpenBlockchainAction failed: %s", err)
		}

		genesisBlock, err := readLatestBlock(db, testBlocksBucketName)
		if err != nil {
			t.Fatalf("error: %s", err)
		}

		if len(genesisBlock.PreviousHash) != 0 {
			t.Fatalf("Genesis block has bad PreviousHash, expected '%x', but was '%x'", []byte{}, genesisBlock.PreviousHash)
		}
		if genesisBlock.Index != 0 {
			t.Fatalf("Genesis block has bad Index, expected '%s', but was '%s'", 0, genesisBlock.Index)
		}
		if !cmp.Equal(head, genesisBlock) {
			t.Fatalf("Resulting block does not equal the latest block: %s", cmp.Diff(head, genesisBlock))
		}
	})

	err = os.Remove(testDbFileName)
	if err != nil {
		t.Fatalf("deleting test db file: %s", err)
	}
}
