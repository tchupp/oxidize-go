package account

import "github.com/tclchiam/oxidize-go/identity"

type Update interface {
	Address() *identity.Address
	Apply(account *Account) *Account
}

type spendUpdate struct {
	address *identity.Address
	amount  uint64
}

func (u *spendUpdate) Address() *identity.Address {
	return u.address
}

func (u *spendUpdate) Apply(account *Account) *Account {
	return &Account{
		address:      account.address,
		spendable:    account.spendable - u.amount,
		transactions: account.transactions,
	}
}

type receiveUpdate struct {
	address *identity.Address
	amount  uint64
}

func (u *receiveUpdate) Address() *identity.Address {
	return u.address
}

func (u *receiveUpdate) Apply(account *Account) *Account {
	return &Account{
		address:      account.address,
		spendable:    account.spendable + u.amount,
		transactions: account.transactions,
	}
}
