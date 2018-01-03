package bolt_impl

import (
	"testing"
	"bytes"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/google/go-cmp/cmp"
	"github.com/tclchiam/block_n_go/tx"
	"github.com/tclchiam/block_n_go/blockchain"
	"github.com/tclchiam/block_n_go/wallet"
	"github.com/tclchiam/block_n_go/chainhash"
)

func TestBlockchainRepository_SaveBlock(t *testing.T) {
	address := wallet.NewWallet().GetAddress()
	const testBlockchainName = "test"

	repository, err := NewRepository(testBlockchainName)
	if err != nil {
		t.Fatalf("creating blockchain repository: %s", err)
	}
	defer closeAndDeleteDB(repository, t)

	transaction := tx.NewCoinbaseTx(address)
	transactions := []*tx.Transaction{transaction}

	previousIndex := 5
	previousHash := chainhash.Hash{}

	err = repository.SaveBlock(blockchain.NewBlock(transactions, previousHash, previousIndex+1))
	if err != nil {
		t.Fatalf("SaveBlock failed: %s", err)
	}

	newBlock, err := readLatestBlock(repository.db, blocksBucketName)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	if bytes.Compare(newBlock.PreviousHash.Slice(), previousHash.Slice()) != 0 {
		t.Fatalf("New block has bad PreviousHash, expected [%s], but was [%s]", previousHash, newBlock.PreviousHash)
	}
	if newBlock.Index != previousIndex+1 {
		t.Fatalf("New block has bad Index, expected [%s], but was [%s]", previousIndex+1, newBlock.Index)
	}
	if !cmp.Equal(newBlock.Transactions, transactions) {
		t.Fatalf("New block has bad Transactions, %s", cmp.Diff(newBlock.Transactions, transactions))
	}
	if !cmp.Equal(newBlock, newBlock) {
		t.Fatalf("Resulting block does not equal the latest block: %s", cmp.Diff(newBlock, newBlock))
	}
}

func readLatestBlock(db *bolt.DB, bucketName []byte) (latestBlock *blockchain.Block, err error) {
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

		latestBlock, err = DeserializeBlock(latestBlockData)
		if err != nil {
			return err
		}

		return nil
	})

	return latestBlock, err
}
