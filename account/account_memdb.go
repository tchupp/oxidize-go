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

func (r *accountRepo) SaveTx(address *identity.Address, tx *Transaction) error {
	if address == nil {
		return nil
	}

	r.lock.Lock()
	defer r.lock.Unlock()
	account, err := r.account(address)
	if err != nil {
		return err
	}

	account.Transactions = append(account.Transactions, tx)
	if tx.spender.IsEqual(address) {
		account.Spendable -= tx.amount
	}
	if tx.receiver.IsEqual(address) {
		account.Spendable += tx.amount
	}

	r.accountStore[address.Serialize()] = account

	return nil
}

func (r *accountRepo) Account(address *identity.Address) (*Account, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.account(address)
}

func (r *accountRepo) account(address *identity.Address) (*Account, error) {
	if a, ok := r.accountStore[address.Serialize()]; ok {
		return a, nil
	} else {
		return &Account{Address: address}, nil
	}
}
