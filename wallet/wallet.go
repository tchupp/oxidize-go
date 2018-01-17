package wallet

import (
	"github.com/tclchiam/block_n_go/crypto"
	"github.com/tclchiam/block_n_go/identity"
)

type Wallet struct {
	PrivateKey *crypto.PrivateKey
	PublicKey  *crypto.PublicKey

	identity *identity.Identity
}

func NewWallet() *Wallet {
	privateKey := crypto.NewP256PrivateKey()
	return newWallet(privateKey)
}

func newWallet(privateKey *crypto.PrivateKey) *Wallet {
	return &Wallet{PrivateKey: privateKey, PublicKey: privateKey.PubKey()}
}

func (w *Wallet) GetAddress() *identity.Identity {
	if w.identity != nil {
		return w.identity
	}

	w.identity = identity.NewIdentity(w.PrivateKey)
	return w.identity
}

func (w *Wallet) String() string {
	return w.GetAddress().String()
}
