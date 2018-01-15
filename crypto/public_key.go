package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

type PublicKey ecdsa.PublicKey

func (p *PublicKey) ToECDSA() *ecdsa.PublicKey {
	return (*ecdsa.PublicKey)(p)
}

func (p *PublicKey) IsEqual(otherPubKey *PublicKey) bool {
	return p.X.Cmp(otherPubKey.X) == 0 && p.Y.Cmp(otherPubKey.Y) == 0
}

func (p *PublicKey) String() string {
	return string(p.Serialize())
}

func (p *PublicKey) Verify(hash []byte, signature *Signature) bool {
	return ecdsa.Verify(p.ToECDSA(), hash, signature.R, signature.S)
}

func (p *PublicKey) Serialize() []byte {
	return append(p.X.Bytes(), p.Y.Bytes()...)
}

func (p *PublicKey) Hash() []byte {
	publicSHA256 := sha256.Sum256(p.Serialize())

	hashImpl := ripemd160.New()
	hashImpl.Write(publicSHA256[:])
	return hashImpl.Sum(nil)
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
