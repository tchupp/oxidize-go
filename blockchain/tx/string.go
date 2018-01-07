package tx

import (
	"fmt"
	"strings"
	"strconv"
)

func (tx *Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %s:", tx.ID))
	lines = append(lines, fmt.Sprintf("     Is Coinbase: %s", strconv.FormatBool(tx.IsCoinbase())))

	for _, input := range tx.TxInputs {
		lines = append(lines, input.String())
	}

	for _, output := range tx.TxOutputs {
		lines = append(lines, output.String())
	}

	return strings.Join(lines, "\n")
}

func (input *UnsignedInput) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("     UnsignedInput:"))
	lines = append(lines, fmt.Sprintf("       TransactionId: %s", input.OutputReference.ID))
	lines = append(lines, fmt.Sprintf("       OutputIndex:   %d", input.OutputReference.Output.Index))
	lines = append(lines, fmt.Sprintf("       PublicKey:     %x", input.PublicKey))
	return strings.Join(lines, "\n")
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

func (output *Output) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("     Output:"))
	lines = append(lines, fmt.Sprintf("       Index:         %d", output.Index))
	lines = append(lines, fmt.Sprintf("       Value:         %d", output.Value))
	lines = append(lines, fmt.Sprintf("       PublicKeyHash: %x", output.PublicKeyHash))

	return strings.Join(lines, "\n")
}