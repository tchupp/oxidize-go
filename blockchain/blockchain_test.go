package blockchain

import (
	"fmt"
	"strings"
	"testing"

	"github.com/boltdb/bolt"
)

func TestBlockchain_Workflow(t *testing.T) {
	const owner = "Theo"
	const actor1 = "Marika"
	const actor2 = "Ivan"
	const actor3 = "Nick"
	const actor4 = "George"

	t.Run("Sending: expense < balance", func(t *testing.T) {
		bc, err := Open("test1", owner)
		if err != nil {
			t.Fatalf("failed to open test blockchain: %s", err)
		}
		defer bc.Delete()

		bc, err = bc.Send(owner, actor1, 3)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, bc, owner, 7)
		verifyBalance(t, bc, actor1, 3)
	})

	t.Run("Sending: expense == balance", func(t *testing.T) {
		bc, err := Open("test2", owner)
		if err != nil {
			t.Fatalf("failed to open test blockchain: %s", err)
		}
		defer bc.Delete()

		bc, err = bc.Send(owner, actor1, 10)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, bc, owner, 0)
		verifyBalance(t, bc, actor1, 10)
	})

	t.Run("Sending: expense > balance", func(t *testing.T) {
		bc, err := Open("test3", owner)
		if err != nil {
			t.Fatalf("failed to open test blockchain: %s", err)
		}
		defer bc.Delete()

		bc, err = bc.Send(owner, actor1, 13)
		if err == nil {
			t.Fatalf("expected error")
		}

		expectedMessage := "account 'Theo' does not have enough to send '13', due to balance '10'"
		if !strings.Contains(err.Error(), expectedMessage) {
			t.Fatalf("Expected string to contain: \"%s\", was '%s'", expectedMessage, err.Error())
		}

		verifyBalance(t, bc, owner, 10)
		verifyBalance(t, bc, actor1, 0)
	})

	t.Run("Sending: many", func(t *testing.T) {
		bc, err := Open("test4", owner)
		if err != nil {
			t.Fatalf("failed to open test blockchain: %s", err)
		}
		defer bc.Delete()

		bc, err = bc.Send(owner, actor1, 1)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, bc, owner, 9)
		verifyBalance(t, bc, actor1, 1)
		verifyBalance(t, bc, actor2, 0)
		verifyBalance(t, bc, actor3, 0)
		verifyBalance(t, bc, actor4, 0)

		bc, err = bc.Send(owner, actor2, 1)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, bc, owner, 8)
		verifyBalance(t, bc, actor1, 1)
		verifyBalance(t, bc, actor1, 1)
		verifyBalance(t, bc, actor3, 0)
		verifyBalance(t, bc, actor4, 0)

		bc, err = bc.Send(owner, actor3, 1)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, bc, owner, 7)
		verifyBalance(t, bc, actor1, 1)
		verifyBalance(t, bc, actor1, 1)
		verifyBalance(t, bc, actor3, 1)
		verifyBalance(t, bc, actor4, 0)

		bc, err = bc.Send(owner, actor4, 1)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, bc, owner, 6)
		verifyBalance(t, bc, actor1, 1)
		verifyBalance(t, bc, actor1, 1)
		verifyBalance(t, bc, actor3, 1)
		verifyBalance(t, bc, actor4, 1)
	})
}

func verifyBalance(t *testing.T, bc *Blockchain, address string, expectedBalance int) {
	balance, err := bc.ReadBalance(address)

	if err != nil {
		t.Fatalf("reading balance for '%s' %s", address, err)
	}
	if balance != expectedBalance {
		t.Fatalf("expected balance for '%s' to be [%d], was: [%d]", address, expectedBalance, balance)
	}
}

func readLatestBlock(db *bolt.DB, bucketName []byte) (latestBlock *Block, err error) {
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
