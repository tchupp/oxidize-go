package entity

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/tclchiam/oxidize-go/crypto"
	"github.com/tclchiam/oxidize-go/identity"
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

func (input *SignedInput) SpentBy(spender *identity.Identity) bool {
	publicKeyHash := spender.PublicKeyHash()
	lockingHash := input.PublicKey.Hash()

	return bytes.Compare(lockingHash, publicKeyHash) == 0
}

func EmptySingedInputs() SignedInputs {
	c := make(chan *SignedInput, 0)
	defer close(c)
	return SignedInputs(c)
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
	for input := range inputs {
		res = apply(res, input)
	}
	return res
}

func (inputs SignedInputs) Add(input *SignedInput) SignedInputs {
	c := make(chan *SignedInput, len(inputs)+1)
	defer close(c)

	for i := range inputs {
		c <- i
	}
	c <- input
	return SignedInputs(c)
}

func (inputs SignedInputs) Append(plus SignedInputs) SignedInputs {
	c := make(chan *SignedInput, len(inputs)+len(plus))

	go func() {
		for i := range inputs {
			c <- i
		}
		for i := range plus {
			c <- i
		}
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
