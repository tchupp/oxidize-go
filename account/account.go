package account

import (
	"fmt"

	"github.com/tclchiam/oxidize-go/identity"
)

type Account struct {
	address      *identity.Address
	spendable    uint64
	transactions Transactions
}

func NewAccount(address *identity.Address, spendable uint64, transactions Transactions) *Account {
	return &Account{
		address:      address,
		spendable:    spendable,
		transactions: transactions,
	}
}

func (a *Account) Address() *identity.Address { return a.address }
func (a *Account) Spendable() uint64          { return a.spendable }
func (a *Account) Transactions() Transactions { return a.transactions }

func (a *Account) IsEqual(other *Account) bool {
	if a == nil && other == nil {
		return true
	}
	if a == nil || other == nil {
		return false
	}
	if a == other {
		return true
	}

	if !a.address.IsEqual(other.address) {
		return false
	}
	if a.spendable != other.spendable {
		return false
	}
	return true
}

func (a *Account) String() string {
	return fmt.Sprintf("{address: '%s', spendable: %d}", a.Address(), a.Spendable())
}
