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
	blockRepository, err := NewBlockRepository(testBlockchainName, encoding.NewBlockGobEncoder())
	if err != nil {
		t.Fatalf("creating block repository: %s", err)
	}
	defer closeAndDeleteDB(blockRepository.(*blockBoltRepository), t)
	blockRepository.Close()

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

func closeAndDeleteDB(repository *blockBoltRepository, t *testing.T) {
	if err := repository.Close(); err != nil {
		t.Fatalf("closing repository: %s", err)
	}

	if err := DeleteBlockchain(repository.name); err != nil {
		t.Fatalf("deleting test db file: %s", err)
	}
}
