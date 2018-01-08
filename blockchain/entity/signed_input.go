package entity

import (
	"bytes"

	"github.com/tclchiam/block_n_go/wallet"
	"github.com/tclchiam/block_n_go/crypto"
	"fmt"
	"strings"
)

type SignedInput struct {
	OutputReference *OutputReference
	Signature       *crypto.Signature
	PublicKey       *crypto.PublicKey
}

func NewSignedInput(input *UnsignedInput, signature *crypto.Signature) *SignedInput {
	return &SignedInput{
		OutputReference: input.OutputReference,
		Signature:       signature,
		PublicKey:       input.PublicKey,
	}
}

func (input *SignedInput) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("     SignedInput:"))
	lines = append(lines, fmt.Sprintf("       TransactionId: %s", input.OutputReference.ID))
	lines = append(lines, fmt.Sprintf("       OutputIndex:   %d", input.OutputReference.Output.Index))
	lines = append(lines, fmt.Sprintf("       PublicKey:     %x", input.PublicKey))
	lines = append(lines, fmt.Sprintf("       Signature:     %x", input.Signature))
	return strings.Join(lines, "\n")
}

func (input *SignedInput) SpentBy(address string) bool {
	publicKeyHash, err := wallet.AddressToPublicKeyHash(address)
	if err != nil {
		return false
	}

	lockingHash := wallet.HashPubKey(input.PublicKey)

	return bytes.Compare(lockingHash, publicKeyHash) == 0
}

func NewSignedInputs(inputs []*SignedInput) SignedInputs {
	c := make(chan *SignedInput, len(inputs))
	go func() {
		for _, input := range inputs {
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
