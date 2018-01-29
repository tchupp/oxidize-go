package blockchain_test

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/tclchiam/block_n_go/blockchain"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/identity"
	"github.com/tclchiam/block_n_go/mining"
	"github.com/tclchiam/block_n_go/mining/proofofwork"
	"github.com/tclchiam/block_n_go/storage/memdb"
)

func TestBlockchain_Workflow(t *testing.T) {
	owner := identity.RandomIdentity()
	actor1 := identity.RandomIdentity()
	actor2 := identity.RandomIdentity()
	actor3 := identity.RandomIdentity()
	actor4 := identity.RandomIdentity()

	t.Run("Sending: expense < balance", func(t *testing.T) {
		bc := setupBlockchain(t, owner)

		err := bc.Send(owner, actor1, owner, 3)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, bc, owner, 17)
		verifyBalance(t, bc, actor1, 3)
	})

	t.Run("Sending: expense == balance", func(t *testing.T) {
		bc := setupBlockchain(t, owner)

		err := bc.Send(owner, actor1, owner, 10)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, bc, owner, 10)
		verifyBalance(t, bc, actor1, 10)
	})

	t.Run("Sending: expense > balance", func(t *testing.T) {
		bc := setupBlockchain(t, owner)

		err := bc.Send(owner, actor1, owner, 13)
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
		bc := setupBlockchain(t, owner)

		err := bc.Send(owner, actor1, owner, 1)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, bc, owner, 19)
		verifyBalance(t, bc, actor1, 1)
		verifyBalance(t, bc, actor2, 0)
		verifyBalance(t, bc, actor3, 0)
		verifyBalance(t, bc, actor4, 0)

		err = bc.Send(owner, actor2, owner, 1)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, bc, owner, 28)
		verifyBalance(t, bc, actor1, 1)
		verifyBalance(t, bc, actor1, 1)
		verifyBalance(t, bc, actor3, 0)
		verifyBalance(t, bc, actor4, 0)

		err = bc.Send(owner, actor3, owner, 1)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, bc, owner, 37)
		verifyBalance(t, bc, actor1, 1)
		verifyBalance(t, bc, actor1, 1)
		verifyBalance(t, bc, actor3, 1)
		verifyBalance(t, bc, actor4, 0)

		err = bc.Send(owner, actor4, owner, 1)
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
	balance, err := bc.ReadBalance(spender)

	if err != nil {
		t.Fatalf("reading balance for '%x' %s", spender, err)
	}
	if balance != expectedBalance {
		t.Fatalf("expected balance for '%x' to be [%d], was: [%d]", spender, expectedBalance, balance)
	}
}

func setupBlockchain(t *testing.T, owner *identity.Identity) blockchain.Blockchain {
	blockRepository := memdb.NewBlockRepository()
	miner := proofofwork.NewDefaultMiner(owner)

	genesisBlock := buildGenesisBlock(miner)
	if err := blockRepository.SaveBlock(genesisBlock); err != nil {
		t.Fatalf("saving genesis block: %s", err)
	}

	bc, err := blockchain.Open(blockRepository, memdb.NewHeaderRepository(), miner)
	if err != nil {
		t.Fatalf("failed to open test blockchain: %s", err)
	}

	return bc
}

func buildGenesisBlock(miner mining.Miner) *entity.Block {
	header := entity.NewBlockHeader(math.MaxUint64, nil, nil, 0, 0, &entity.EmptyHash)
	parent := entity.NewBlock(header, entity.Transactions{})

	return miner.MineBlock(parent, entity.Transactions{})
}
