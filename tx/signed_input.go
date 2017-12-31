package tx

import (
	"bytes"

	"github.com/tclchiam/block_n_go/wallet"
)

type SignedInput struct {
	OutputReference OutputReference
	Signature       []byte
	PublicKey       []byte
}

func NewSignedInput(input *UnsignedInput, signature []byte) *SignedInput {
	return &SignedInput{
		OutputReference: input.OutputReference,
		Signature:       signature,
		PublicKey:       input.PublicKey,
	}
}

func (input *SignedInput) SpentBy(address string) bool {
	publicKeyHash, err := wallet.AddressToPublicKeyHash(address)
	if err != nil {
		return false
	}

	lockingHash := wallet.HashPubKey(input.PublicKey)

	return bytes.Compare(lockingHash, publicKeyHash) == 0
}
