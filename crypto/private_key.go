package crypto

import (
	log "github.com/sirupsen/logrus"

	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

type PrivateKey ecdsa.PrivateKey

func NewPrivateKey(curve elliptic.Curve) *PrivateKey {
	key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	return (*PrivateKey)(key)
}

func NewP256PrivateKey() *PrivateKey {
	return NewPrivateKey(elliptic.P256())
}

func FromPrivateKey(key *ecdsa.PrivateKey) *PrivateKey {
	return (*PrivateKey)(key)
}

func (p *PrivateKey) ToECDSA() *ecdsa.PrivateKey {
	return (*ecdsa.PrivateKey)(p)
}

func (p *PrivateKey) PubKey() *PublicKey {
	return (*PublicKey)(&p.PublicKey)
}

func (p *PrivateKey) Sign(hash []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, p.ToECDSA(), hash)
	if err != nil {
		return nil, err
	}
	return &Signature{R: r, S: s}, nil
}
