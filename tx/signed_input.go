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

func newSignedInput(input *UnsignedInput, signature []byte) *SignedInput {
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

func (tx *Transaction) Inputs() SignedInputs {
	c := make(chan *SignedInput, len(tx.TxInputs))
	go func() {
		for _, input := range tx.TxInputs {
			c <- input
		}
		close(c)
	}()
	return SignedInputs(c)
}

type SignedInputs <-chan *SignedInput

func (inputs SignedInputs) Filter(predicate func(input *SignedInput) bool) SignedInputs {
	c := make(chan *SignedInput)

	go func() {
		for input := range inputs {
			if predicate(input) {
				c <- input
			}
		}
		close(c)
	}()
	return SignedInputs(c)
}

func (inputs SignedInputs) Reduce(res interface{}, apply func(res interface{}, input *SignedInput) interface{}) interface{} {
	c := make(chan interface{})

	go func() {
		for input := range inputs {
			res = apply(res, input)
		}
		c <- res
	}()
	return <-c
}

func (inputs SignedInputs) Add(input *SignedInput) SignedInputs {
	c := make(chan *SignedInput)

	go func() {
		for i := range inputs {
			c <- i
		}
		c <- input
		close(c)
	}()
	return SignedInputs(c)
}

func (inputs SignedInputs) ToSlice() []*SignedInput {
	slice := make([]*SignedInput, 0)
	for i := range inputs {
		slice = append(slice, i)
	}
	return slice
}