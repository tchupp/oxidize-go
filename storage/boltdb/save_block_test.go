package boltdb

import (
	"testing"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/google/go-cmp/cmp"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/blockchain/entity/encoding"
	"github.com/tclchiam/block_n_go/identity"
	"github.com/tclchiam/block_n_go/mining"
)

func TestBlockRepository_SaveBlock(t *testing.T) {
	randomIdentity := identity.RandomIdentity()
	const testBlockchainName = "test"

	blockEncoder := encoding.NewBlockGobEncoder()
	blockRepository, err := NewBlockRepository(testBlockchainName, blockEncoder)
	if err != nil {
		t.Fatalf("creating block repository: %s", err)
	}
	defer deleteDB(testBlockchainName, t)

	transactions := entity.Transactions{entity.NewCoinbaseTx(randomIdentity, encoding.TransactionProtoEncoder())}

	const previousIndex = 5
	previousHash := entity.NewHashOrPanic("0000f65fe866ab6f810b13a5d864f96cb16ad22e2e931b861f80d870f2e32df7")
	hash := entity.NewHashOrPanic("00007eaa535b8894e8815f57d317c3bb14ab598417fe4ddd8d37d65c189f85fe")

	blockToSave := entity.NewBlock(
		entity.NewBlockHeader(previousIndex+1, previousHash, mining.CalculateTransactionsHash(transactions), 18920304, 38385, hash),
		transactions,
	)

	err = blockRepository.SaveBlock(blockToSave)
	if err != nil {
		t.Fatalf("SaveBlock failed: %s", err)
	}

	db, err := openDB(testBlockchainName)
	if err != nil {
		t.Fatalf("opening database: %s", err)
	}

	newBlock, err := readLatestBlock(db, blocksBucketName, blockEncoder)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	if !blockToSave.IsEqual(newBlock) {
		t.Fatalf("Resulting block does not equal the latest block: %s", cmp.Diff(blockToSave, newBlock))
	}
}

func readLatestBlock(db *bolt.DB, bucketName []byte, encoder entity.BlockEncoder) (latestBlock *entity.Block, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("no block with name %s exists", bucketName)
		}

		latestBlockHash := bucket.Get([]byte("l"))
		if latestBlockHash == nil {
			return fmt.Errorf("could not find latest block hash")
		}

		latestBlockData := bucket.Get(latestBlockHash)
		if latestBlockData == nil || len(latestBlockData) == 0 {
			return fmt.Errorf("latest block data is empty: '%s'", latestBlockData)
		}

		latestBlock, err = encoder.DecodeBlock(latestBlockData)
		if err != nil {
			return err
		}

		return nil
	})

	return latestBlock, err
}
