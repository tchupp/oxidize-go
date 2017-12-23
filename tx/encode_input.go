package tx

import "fmt"

func (input *Input) string(id int) []string {
	var lines []string
	lines = append(lines, fmt.Sprintf("     Input:"))
	lines = append(lines, fmt.Sprintf("       Id:            %x", id))
	lines = append(lines, fmt.Sprintf("       TransactionId: %x", input.OutputTransactionId))
	lines = append(lines, fmt.Sprintf("       OutputId:      %d", input.OutputId))
	lines = append(lines, fmt.Sprintf("       Signature:     %x", input.Signature))
	lines = append(lines, fmt.Sprintf("       PublicKey:     %x", input.PublicKey))
	return lines
}
