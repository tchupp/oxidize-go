package wallet

import (
	"github.com/tclchiam/block_n_go/crypto"
	"github.com/tclchiam/block_n_go/identity"
)

type Wallet struct {
	PrivateKey *crypto.PrivateKey
	PublicKey  *crypto.PublicKey

	address *identity.Address
}

func NewWallet() *Wallet {
	privateKey := crypto.NewP256PrivateKey()
	return newWallet(privateKey)
}

func newWallet(privateKey *crypto.PrivateKey) *Wallet {
	return &Wallet{PrivateKey: privateKey, PublicKey: privateKey.PubKey()}
}

func (w *Wallet) GetAddress() *identity.Address {
	if w.address != nil {
		return w.address
	}

	w.address = identity.NewAddress(w.PrivateKey)
	return w.address
}

func (w *Wallet) String() string {
	return w.GetAddress().String()
}
