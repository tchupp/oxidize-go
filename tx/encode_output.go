package tx

import "fmt"

func (output *Output) string() []string {
	var lines []string

	lines = append(lines, fmt.Sprintf("     Output:"))
	lines = append(lines, fmt.Sprintf("       Id:            %d", output.Id))
	lines = append(lines, fmt.Sprintf("       Value:         %d", output.Value))
	lines = append(lines, fmt.Sprintf("       PublicKeyHash: %x", output.PublicKeyHash))

	return lines
}
