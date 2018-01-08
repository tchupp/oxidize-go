package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
)

type PublicKey ecdsa.PublicKey

func (p *PublicKey) ToECDSA() *ecdsa.PublicKey {
	return (*ecdsa.PublicKey)(p)
}

func (p *PublicKey) IsEqual(otherPubKey *PublicKey) bool {
	return p.X.Cmp(otherPubKey.X) == 0 && p.Y.Cmp(otherPubKey.Y) == 0
}

func (p *PublicKey) Serialize() []byte {
	return append(p.X.Bytes(), p.Y.Bytes()...)
}

func DeserializePublicKey(input []byte) (*PublicKey, error) {
	inputLen := len(input)

	publicKey := &PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int).SetBytes(input[:(inputLen / 2)]),
		Y:     new(big.Int).SetBytes(input[(inputLen / 2):]),
	}

	return publicKey, nil
}
