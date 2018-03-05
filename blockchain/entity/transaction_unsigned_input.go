package entity

import (
	"fmt"
	"strings"

	"github.com/tclchiam/oxidize-go/crypto"
)

type UnsignedInput struct {
	OutputReference *OutputReference
	PublicKey       *crypto.PublicKey
}

func NewUnsignedInput(outputTransactionId *Hash, output *Output, spenderPublicKey *crypto.PublicKey) *UnsignedInput {
	reference := &OutputReference{ID: outputTransactionId, Output: output}

	return &UnsignedInput{
		OutputReference: reference,
		PublicKey:       spenderPublicKey,
	}
}

func (input *UnsignedInput) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("     UnsignedInput:"))
	lines = append(lines, fmt.Sprintf("       TransactionId: %s", input.OutputReference.ID))
	lines = append(lines, fmt.Sprintf("       OutputIndex:   %d", input.OutputReference.Output.Index))
	lines = append(lines, fmt.Sprintf("       PublicKey:     %x", input.PublicKey))
	return strings.Join(lines, "\n")
}
