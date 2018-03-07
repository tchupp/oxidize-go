package account

import (
	"sync"

	"github.com/tclchiam/oxidize-go/identity"
)

type accountRepo struct {
	lock         sync.RWMutex
	accountStore map[string]*Account
}

func NewAccountRepository() *accountRepo {
	return &accountRepo{accountStore: make(map[string]*Account)}
}

func (r *accountRepo) ProcessUpdates(updates []Update) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, update := range updates {
		account := r.account(update.Address())
		updatedAccount := update.Apply(account)

		r.accountStore[update.Address().Serialize()] = updatedAccount
	}

	return nil
}

func (r *accountRepo) Account(address *identity.Address) (*Account, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.account(address), nil
}

func (r *accountRepo) account(address *identity.Address) *Account {
	if a, ok := r.accountStore[address.Serialize()]; ok {
		return a
	} else {
		return &Account{address: address}
	}
}
