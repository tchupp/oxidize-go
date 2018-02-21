package account

import (
	"fmt"

	"github.com/tclchiam/oxidize-go/identity"
)

type Account struct {
	Address      *identity.Address
	Spendable    uint64
	Transactions Transactions
}

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

	if !a.Address.IsEqual(other.Address) {
		return false
	}
	if a.Spendable != other.Spendable {
		return false
	}
	return true
}

func (a *Account) String() string {
	return fmt.Sprintf("{address: '%s', spendable: %d}", a.Address, a.Spendable)
}
