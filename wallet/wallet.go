package wallet

import (
	"github.com/tclchiam/oxidize-go/account"
	"github.com/tclchiam/oxidize-go/identity"
	"github.com/tclchiam/oxidize-go/wallet/rpc"
)

type Wallet interface {
	Identities() (identity.Identities, error)
	NewIdentity() (*identity.Identity, error)
	Account() ([]*account.Account, error)
}

type wallet struct {
	store  *KeyStore
	client rpc.WalletClient
}

func NewWallet(store *KeyStore, client rpc.WalletClient) Wallet {
	return &wallet{store: store, client: client}
}

func (w *wallet) Identities() (identity.Identities, error) {
	return w.store.Identities()
}

func (w *wallet) NewIdentity() (*identity.Identity, error) {
	newIdentity := identity.RandomIdentity()
	err := w.store.SaveIdentity(newIdentity)
	if err != nil {
		return nil, err
	}

	return newIdentity, nil
}

func (w *wallet) Account() ([]*account.Account, error) {
	identities, err := w.store.Identities()
	if err != nil {
		return nil, err
	}

	var addrs []*identity.Address
	for _, id := range identities {
		addrs = append(addrs, id.Address())
	}

	return w.client.Account(addrs)
}
