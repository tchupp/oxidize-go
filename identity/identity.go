package identity

import (
	"github.com/tclchiam/oxidize-go/crypto"
)

type Identity struct {
	publicKey  *crypto.PublicKey
	privateKey *crypto.PrivateKey

	address *Address
}

func RandomIdentity() *Identity {
	privateKey := crypto.NewP256PrivateKey()
	return NewIdentity(privateKey)
}

func NewIdentity(privateKey *crypto.PrivateKey) *Identity {
	return &Identity{
		publicKey:  privateKey.PubKey(),
		privateKey: privateKey,
	}
}

func (a *Identity) Address() *Address {
	if a.address != nil {
		return a.address
	}

	a.address = FromPublicKey(a.publicKey)
	return a.address
}

func (a *Identity) Sign(data []byte) (*crypto.Signature, error) {
	return a.privateKey.Sign(data)
}

func (a *Identity) String() string                 { return a.Address().String() }
func (a *Identity) PrivateKey() *crypto.PrivateKey { return a.privateKey }
func (a *Identity) PublicKey() *crypto.PublicKey   { return a.publicKey }
func (a *Identity) IsEqual(other *Identity) bool   { return a.Address().IsEqual(other.Address()) }

type Identities []*Identity

func (i Identities) FirstOrNil() *Identity {
	if len(i) > 0 {
		return i[0]
	}
	return nil
}
