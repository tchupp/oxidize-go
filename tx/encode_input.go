package tx

import "fmt"

func (input *UnsignedInput) string(id int) []string {
	var lines []string
	lines = append(lines, fmt.Sprintf("     UnsignedInput:"))
	lines = append(lines, fmt.Sprintf("       Id:            %x", id))
	lines = append(lines, fmt.Sprintf("       TransactionId: %x", input.OutputTransactionId))
	lines = append(lines, fmt.Sprintf("       OutputId:      %d", input.OutputId))
	lines = append(lines, fmt.Sprintf("       PublicKey:     %x", input.PublicKey))
	return lines
}
