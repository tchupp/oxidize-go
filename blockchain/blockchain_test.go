package blockchain_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/blockchain/engine/mining/proofofwork"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/encoding"
	"github.com/tclchiam/oxidize-go/identity"
	"github.com/tclchiam/oxidize-go/storage/boltdb"
)

func TestBlockchain_Workflow(t *testing.T) {
	owner := identity.RandomIdentity()
	actor1 := identity.RandomIdentity()
	actor2 := identity.RandomIdentity()
	actor3 := identity.RandomIdentity()
	actor4 := identity.RandomIdentity()

	t.Run("Sending: expense < balance", func(t *testing.T) {
		const name = "test1"
		bc := setupBlockchain(t, name, owner)
		defer boltdb.DeleteBlockchain(name)

		err := bc.Send(owner, actor1.Address(), 3)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, bc, owner, 17)
		verifyBalance(t, bc, actor1, 3)
	})

	t.Run("Sending: expense == balance", func(t *testing.T) {
		const name = "test2"
		bc := setupBlockchain(t, name, owner)
		defer boltdb.DeleteBlockchain(name)

		err := bc.Send(owner, actor1.Address(), 10)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, bc, owner, 10)
		verifyBalance(t, bc, actor1, 10)
	})

	t.Run("Sending: expense > balance", func(t *testing.T) {
		const name = "test3"
		bc := setupBlockchain(t, name, owner)
		defer boltdb.DeleteBlockchain(name)

		err := bc.Send(owner, actor1.Address(), 13)
		if err == nil {
			t.Fatalf("expected error")
		}

		expectedMessage := fmt.Sprintf("account '%s' does not have enough to send '13', due to balance '10'", owner)
		if !strings.Contains(err.Error(), expectedMessage) {
			t.Fatalf("Expected string to contain: \"%s\", was '%s'", expectedMessage, err.Error())
		}

		verifyBalance(t, bc, owner, 10)
		verifyBalance(t, bc, actor1, 0)
	})

	t.Run("Sending: many", func(t *testing.T) {
		const name = "test4"
		bc := setupBlockchain(t, name, owner)
		defer boltdb.DeleteBlockchain(name)

		err := bc.Send(owner, actor1.Address(), 1)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, bc, owner, 19)
		verifyBalance(t, bc, actor1, 1)
		verifyBalance(t, bc, actor2, 0)
		verifyBalance(t, bc, actor3, 0)
		verifyBalance(t, bc, actor4, 0)

		err = bc.Send(owner, actor2.Address(), 1)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, bc, owner, 28)
		verifyBalance(t, bc, actor1, 1)
		verifyBalance(t, bc, actor1, 1)
		verifyBalance(t, bc, actor3, 0)
		verifyBalance(t, bc, actor4, 0)

		err = bc.Send(owner, actor3.Address(), 1)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, bc, owner, 37)
		verifyBalance(t, bc, actor1, 1)
		verifyBalance(t, bc, actor1, 1)
		verifyBalance(t, bc, actor3, 1)
		verifyBalance(t, bc, actor4, 0)

		err = bc.Send(owner, actor4.Address(), 1)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, bc, owner, 46)
		verifyBalance(t, bc, actor1, 1)
		verifyBalance(t, bc, actor1, 1)
		verifyBalance(t, bc, actor3, 1)
		verifyBalance(t, bc, actor4, 1)
	})
}

func verifyBalance(t *testing.T, bc blockchain.Blockchain, spender *identity.Identity, expectedBalance uint32) {
	balance, err := bc.Balance(spender.Address())

	if err != nil {
		t.Fatalf("reading balance for '%s': %s", spender, err)
	}
	if balance != expectedBalance {
		t.Fatalf("expected balance for '%s' to be [%d], was: [%d]", spender, expectedBalance, balance)
	}
}

func setupBlockchain(t *testing.T, name string, owner *identity.Identity) blockchain.Blockchain {
	repository := boltdb.Builder(name, encoding.BlockProtoEncoder()).
		WithCache().
		WithLogger().
		Build()
	miner := proofofwork.NewDefaultMiner(owner.Address())

	genesisBlock := miner.MineBlock(&entity.GenesisParentHeader, entity.Transactions{})
	if err := repository.SaveBlock(genesisBlock); err != nil {
		t.Fatalf("saving genesis block: %s", err)
	}

	bc, err := blockchain.Open(repository, miner)
	if err != nil {
		t.Fatalf("failed to open test blockchain: %s", err)
	}

	return bc
}
