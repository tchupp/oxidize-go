package main

import (
	"strings"
	"testing"

	"github.com/tclchiam/block_n_go/bolt_impl"
	"github.com/tclchiam/block_n_go/blockchain"
	"github.com/tclchiam/block_n_go/wallet"
	"fmt"
)

func TestBlockchain_Workflow(t *testing.T) {
	owner := wallet.NewWallet()
	actor1 := wallet.NewWallet()
	actor2 := wallet.NewWallet()
	actor3 := wallet.NewWallet()
	actor4 := wallet.NewWallet()

	t.Run("Sending: expense < balance", func(t *testing.T) {
		const name = "test1"

		repository, err := bolt_impl.NewRepository(name)
		if err != nil {
			t.Fatalf("failed to create blockchain repository: %s", err)
		}

		bc, err := blockchain.Open(repository, owner.GetAddress())
		if err != nil {
			t.Fatalf("failed to open test blockchain: %s", err)
		}
		defer bolt_impl.DeleteBlockchain(name)

		err = bc.Send(owner, actor1, owner, 3)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, bc, owner, 17)
		verifyBalance(t, bc, actor1, 3)
	})

	t.Run("Sending: expense == balance", func(t *testing.T) {
		const name = "test2"

		repository, err := bolt_impl.NewRepository(name)
		if err != nil {
			t.Fatalf("failed to create blockchain repository: %s", err)
		}

		bc, err := blockchain.Open(repository, owner.GetAddress())
		if err != nil {
			t.Fatalf("failed to open test blockchain: %s", err)
		}
		defer bolt_impl.DeleteBlockchain(name)

		err = bc.Send(owner, actor1, owner, 10)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, bc, owner, 10)
		verifyBalance(t, bc, actor1, 10)
	})

	t.Run("Sending: expense > balance", func(t *testing.T) {
		const name = "test3"

		repository, err := bolt_impl.NewRepository(name)
		if err != nil {
			t.Fatalf("failed to create blockchain repository: %s", err)
		}

		bc, err := blockchain.Open(repository, owner.GetAddress())
		if err != nil {
			t.Fatalf("failed to open test blockchain: %s", err)
		}
		defer bolt_impl.DeleteBlockchain(name)

		err = bc.Send(owner, actor1, owner, 13)
		if err == nil {
			t.Fatalf("expected error")
		}

		expectedMessage := fmt.Sprintf("account '%s' does not have enough to send '13', due to balance '10'", owner.GetAddress())
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

		bc, err := blockchain.Open(repository, owner.GetAddress())
		if err != nil {
			t.Fatalf("failed to open test blockchain: %s", err)
		}
		defer bolt_impl.DeleteBlockchain(name)

		err = bc.Send(owner, actor1, owner, 1)
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

func verifyBalance(t *testing.T, bc *blockchain.Blockchain, wallet *wallet.Wallet, expectedBalance int) {
	address := wallet.GetAddress()
	balance, err := bc.ReadBalance(address)

	if err != nil {
		t.Fatalf("reading balance for '%x' %s", address, err)
	}
	if balance != uint(expectedBalance) {
		t.Fatalf("expected balance for '%x' to be [%d], was: [%d]", address, expectedBalance, balance)
	}
}
