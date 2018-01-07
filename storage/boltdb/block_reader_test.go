package boltdb

import (
	"testing"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/tclchiam/block_n_go/storage"
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
	reader, err := NewReader(testBlockchainName)
	if err != nil {
		t.Fatalf("creating block reader: %s", err)
	}
	defer closeAndDeleteDB(reader.(*blockReader), t)
	reader.Close()

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

func closeAndDeleteDB(reader *blockReader, t *testing.T) {
	if err := reader.Close(); err != nil {
		t.Fatalf("closing reader: %s", err)
	}

	if err := DeleteBlockchain(reader.name); err != nil {
		t.Fatalf("deleting test db file: %s", err)
	}
}
