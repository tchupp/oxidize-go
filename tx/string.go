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

	for id, input := range tx.TxInputs {
		lines = append(lines, input.string(id)...)
	}

	for _, output := range tx.TxOutputs {
		lines = append(lines, output.string()...)
	}

	return strings.Join(lines, "\n")
}
