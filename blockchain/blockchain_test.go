package blockchain_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tclchiam/oxidize-go/account"
	"github.com/tclchiam/oxidize-go/account/testdata"
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
		engine := setupAccountEngine(t, name, owner)
		defer boltdb.DeleteBlockchain(name)

		err := engine.Send(owner, actor1.Address(), 3)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, engine, owner, 17)
		verifyBalance(t, engine, actor1, 3)
	})

	t.Run("Sending: expense == balance", func(t *testing.T) {
		const name = "test2"
		engine := setupAccountEngine(t, name, owner)
		defer boltdb.DeleteBlockchain(name)

		err := engine.Send(owner, actor1.Address(), 10)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, engine, owner, 10)
		verifyBalance(t, engine, actor1, 10)
	})

	t.Run("Sending: expense > balance", func(t *testing.T) {
		const name = "test3"
		engine := setupAccountEngine(t, name, owner)
		defer boltdb.DeleteBlockchain(name)

		err := engine.Send(owner, actor1.Address(), 13)
		if err == nil {
			t.Fatalf("expected error")
		}

		expectedMessage := fmt.Sprintf("account '%s' does not have enough to send '13', due to balance '10'", owner)
		if !strings.Contains(err.Error(), expectedMessage) {
			t.Fatalf("Expected string to contain: \"%s\", was '%s'", expectedMessage, err.Error())
		}

		verifyBalance(t, engine, owner, 10)
		verifyBalance(t, engine, actor1, 0)
	})

	t.Run("Sending: many", func(t *testing.T) {
		const name = "test4"
		engine := setupAccountEngine(t, name, owner)
		defer boltdb.DeleteBlockchain(name)

		err := engine.Send(owner, actor1.Address(), 1)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, engine, owner, 19)
		verifyBalance(t, engine, actor1, 1)
		verifyBalance(t, engine, actor2, 0)
		verifyBalance(t, engine, actor3, 0)
		verifyBalance(t, engine, actor4, 0)

		err = engine.Send(owner, actor2.Address(), 1)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, engine, owner, 28)
		verifyBalance(t, engine, actor1, 1)
		verifyBalance(t, engine, actor1, 1)
		verifyBalance(t, engine, actor3, 0)
		verifyBalance(t, engine, actor4, 0)

		err = engine.Send(owner, actor3.Address(), 1)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, engine, owner, 37)
		verifyBalance(t, engine, actor1, 1)
		verifyBalance(t, engine, actor1, 1)
		verifyBalance(t, engine, actor3, 1)
		verifyBalance(t, engine, actor4, 0)

		err = engine.Send(owner, actor4.Address(), 1)
		if err != nil {
			t.Fatalf("error sending: %s", err)
		}

		verifyBalance(t, engine, owner, 46)
		verifyBalance(t, engine, actor1, 1)
		verifyBalance(t, engine, actor1, 1)
		verifyBalance(t, engine, actor3, 1)
		verifyBalance(t, engine, actor4, 1)
	})
}

func verifyBalance(t *testing.T, engine account.Engine, spender *identity.Identity, expectedBalance uint64) {
	balance, err := engine.Account(spender.Address())

	if err != nil {
		t.Fatalf("reading balance for '%s': %s", spender, err)
	}
	if balance.Spendable() != expectedBalance {
		t.Fatalf("expected balance for '%s' to be [%d], was: [%d]", spender, expectedBalance, balance.Spendable())
	}
}

func setupAccountEngine(t *testing.T, name string, owner *identity.Identity) account.Engine {
	repository := boltdb.ChainBuilder(name).
		WithCache().
		WithLogger().
		Build()

	return testdata.NewAccountEngineBuilder(t).
		WithRepository(repository).
		WithBeneficiary(owner).
		Build().
		AddBalance(owner.Address(), 0)
}
