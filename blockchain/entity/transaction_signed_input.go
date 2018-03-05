package entity

import (
	"fmt"

	"github.com/tclchiam/oxidize-go/crypto"
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
	return input.string("")
}

func (input *SignedInput) string(indent string) string {
	return fmt.Sprintf(
		"%sentity.SignedInput{TransactionId: %s, OutputIndex: %d, PublicKey: %x, Signature: %x}",
		indent,
		input.OutputReference.ID,
		input.OutputReference.Output.Index,
		input.PublicKey,
		input.Signature,
	)
}
