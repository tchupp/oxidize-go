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

func (r *accountRepo) SaveTxs(txs Transactions) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, tx := range txs {
		r.saveTx(tx.spender, tx)
		r.saveTx(tx.receiver, tx)
	}

	return nil
}

func (r *accountRepo) SaveTx(tx *Transaction) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.saveTx(tx.spender, tx)
	r.saveTx(tx.receiver, tx)

	return nil
}

func (r *accountRepo) saveTx(address *identity.Address, tx *Transaction) {
	if address == nil {
		return
	}

	account := r.account(address)

	account.Transactions = account.Transactions.Add(tx)
	if tx.spender.IsEqual(address) {
		account.Spendable -= tx.amount
	}
	if tx.receiver.IsEqual(address) {
		account.Spendable += tx.amount
	}

	r.accountStore[address.Serialize()] = account
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
		return &Account{Address: address}
	}
}
