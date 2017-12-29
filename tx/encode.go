package tx

import (
	"bytes"
	"encoding/gob"
	"log"
	"fmt"
	"strings"
	"strconv"
)

func serialize(tx *Transaction) []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

func (tx *Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))
	lines = append(lines, fmt.Sprintf("     Is Coinbase: %s", strconv.FormatBool(tx.IsCoinbase())))

	for id, input := range tx.TxInputs {
		lines = append(lines, input.string(id)...)
	}

	for _, output := range tx.TxOutputs {
		lines = append(lines, output.string()...)
	}

	return strings.Join(lines, "\n")
}
