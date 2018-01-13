package boltdb

import (
	"testing"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/tclchiam/block_n_go/blockchain/entity/encoding"
)

func TestNewRepository(t *testing.T) {
	const testBlockchainName = "test"

	// Verify starting state
	db, err := openDB(testBlockchainName)
	if err != nil {
		t.Fatalf("opening database: %s", err)
	}
	defer deleteDB(testBlockchainName, t)

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(blocksBucketName)

		if bucket != nil {
			return fmt.Errorf("expected '%s' to be nil", blocksBucketName)
		}

		return nil
	})
	db.Close()
	if err != nil {
		t.Fatalf("Expected no error: %s", err)
	}

	// Execute
	_, err = NewBlockRepository(testBlockchainName, encoding.NewBlockGobEncoder())
	if err != nil {
		t.Fatalf("creating block repository: %s", err)
	}

	// Verify execution
	db, err = openDB(testBlockchainName)
	if err != nil {
		t.Fatalf("opening database: %s", err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(blocksBucketName)

		if bucket == nil {
			return fmt.Errorf("expected '%s' to exist", blocksBucketName)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Expected no error: %s", err)
	}
}

func deleteDB(name string, t *testing.T) {
	if err := DeleteBlockchain(name); err != nil {
		t.Fatalf("deleting test db file: %s", err)
	}
}
