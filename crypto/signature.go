package crypto

import (
	"crypto/ecdsa"
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

func (sig *Signature) Verify(hash []byte, pubKey *PublicKey) bool {
	return ecdsa.Verify(pubKey.ToECDSA(), hash, sig.R, sig.S)
}

func (sig *Signature) IsEqual(otherSig *Signature) bool {
	return sig.R.Cmp(otherSig.R) == 0 && sig.S.Cmp(otherSig.S) == 0
}

func (sig *Signature) String() string {
	return string(sig.Serialize())
}

func (sig *Signature) Serialize() []byte {
	return append(sig.R.Bytes(), sig.S.Bytes()...)
}

func DeserializeSignature(input []byte) (*Signature, error) {
	inputLen := len(input)

	signature := &Signature{
		R: new(big.Int).SetBytes(input[:(inputLen / 2)]),
		S: new(big.Int).SetBytes(input[(inputLen / 2):]),
	}

	return signature, nil
}
