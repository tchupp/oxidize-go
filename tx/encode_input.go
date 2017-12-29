package tx

import "fmt"

func (input *UnsignedInput) string(id int) []string {
	var lines []string
	lines = append(lines, fmt.Sprintf("     UnsignedInput:"))
	lines = append(lines, fmt.Sprintf("       Id:            %x", id))
	lines = append(lines, fmt.Sprintf("       TransactionId: %x", input.OutputReference.ID))
	lines = append(lines, fmt.Sprintf("       OutputIndex:   %d", input.OutputReference.OutputIndex))
	lines = append(lines, fmt.Sprintf("       PublicKey:     %x", input.PublicKey))
	return lines
}

func (input *SignedInput) string(id int) []string {
	var lines []string
	lines = append(lines, fmt.Sprintf("     SignedInput:"))
	lines = append(lines, fmt.Sprintf("       Id:            %x", id))
	lines = append(lines, fmt.Sprintf("       TransactionId: %x", input.OutputReference.ID))
	lines = append(lines, fmt.Sprintf("       OutputIndex:   %d", input.OutputReference.OutputIndex))
	lines = append(lines, fmt.Sprintf("       Signature:     %x", input.Signature))
	lines = append(lines, fmt.Sprintf("       PublicKey:     %x", input.PublicKey))
	return lines
}
