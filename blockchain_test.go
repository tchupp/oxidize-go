package main

import (
	"strings"
	"testing"

	"github.com/tclchiam/block_n_go/bolt_impl"
	"github.com/tclchiam/block_n_go/blockchain"
	"github.com/tclchiam/block_n_go/wallet"
)

func TestBlockchain_Workflow(t *testing.T) {
	owner := wallet.NewWallet().GetAddress()
	actor1 := wallet.NewWallet().GetAddress()
	actor2 := wallet.NewWallet().GetAddress()
	actor3 := wallet.NewWallet().GetAddress()
	actor4 := wallet.NewWallet().GetAddress()

	t.Run("Sending: expense < balance", func(t *testing.T) {
		const name = "test1"

		repository, err := bolt_impl.NewRepository(name)
		if err != nil {
			t.Fatalf("failed to create blockchain repository: %s", err)
		}

		bc, err := blockchain.Open(repository, owner)
		if err != nil {
			t.Fatalf("failed to open test blockchain: %s", err)
		}
		defer bolt_impl.DeleteBlockchain(name)

		bc, err = bc.Send(owner, actor1, 3)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, bc, owner, 7)
		verifyBalance(t, bc, actor1, 3)
	})

	t.Run("Sending: expense == balance", func(t *testing.T) {
		const name = "test2"

		repository, err := bolt_impl.NewRepository(name)
		if err != nil {
			t.Fatalf("failed to create blockchain repository: %s", err)
		}

		bc, err := blockchain.Open(repository, owner)
		if err != nil {
			t.Fatalf("failed to open test blockchain: %s", err)
		}
		defer bolt_impl.DeleteBlockchain(name)

		bc, err = bc.Send(owner, actor1, 10)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, bc, owner, 0)
		verifyBalance(t, bc, actor1, 10)
	})

	t.Run("Sending: expense > balance", func(t *testing.T) {
		const name = "test3"

		repository, err := bolt_impl.NewRepository(name)
		if err != nil {
			t.Fatalf("failed to create blockchain repository: %s", err)
		}

		bc, err := blockchain.Open(repository, owner)
		if err != nil {
			t.Fatalf("failed to open test blockchain: %s", err)
		}
		defer bolt_impl.DeleteBlockchain(name)

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
		const name = "test4"

		repository, err := bolt_impl.NewRepository(name)
		if err != nil {
			t.Fatalf("failed to create blockchain repository: %s", err)
		}

		bc, err := blockchain.Open(repository, owner)
		if err != nil {
			t.Fatalf("failed to open test blockchain: %s", err)
		}
		defer bolt_impl.DeleteBlockchain(name)

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

func verifyBalance(t *testing.T, bc *blockchain.Blockchain, address string, expectedBalance int) {
	balance, err := bc.ReadBalance(address)

	if err != nil {
		t.Fatalf("reading balance for '%s' %s", address, err)
	}
	if balance != expectedBalance {
		t.Fatalf("expected balance for '%s' to be [%d], was: [%d]", address, expectedBalance, balance)
	}
}
