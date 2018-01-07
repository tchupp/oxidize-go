package boltdb

import (
	"testing"
	"fmt"
	"github.com/boltdb/bolt"
	"os"
)

func TestNewRepository(t *testing.T) {
	const testBlockchainName = "test"

	// Verify starting state
	path := fmt.Sprintf(dbFile, testBlockchainName)
	db, err := openDB(path)
	if err != nil {
		t.Fatalf("opening database: %s", err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(blocksBucketName)

		if bucket != nil {
			return fmt.Errorf("expected '%s' to be nil", blocksBucketName)
		}

		return nil
	})
	if err != nil {
		t.Fatalf("Expected no error: %s", err)
	}
	db.Close()

	// Execute
	repository, err := NewRepository(testBlockchainName)
	if err != nil {
		t.Fatalf("creating blockchain repository: %s", err)
	}
	defer closeAndDeleteDB(repository, t)
	repository.Close()

	// Verify execution
	db, err = openDB(path)
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

func closeAndDeleteDB(repository *BlockchainRepository, t *testing.T) {
	repository.Close()

	if err := os.Remove(repository.Path); err != nil {
		t.Fatalf("deleting test db file: %s", err)
	}
}
