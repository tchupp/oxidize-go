package tx

import (
	"bytes"
	"encoding/gob"
	"log"
	"fmt"
	"strings"
	"crypto/sha256"
	"strconv"
)

func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

func (tx *Transaction) Serialize() []byte {
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

func (input *Input) string(id int) []string {
	var lines []string
	lines = append(lines, fmt.Sprintf("     Input %d:", id))
	lines = append(lines, fmt.Sprintf("       TXID:      %x", input.OutputTransactionId))
	lines = append(lines, fmt.Sprintf("       Out:       %d", input.OutputId))
	//lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
	lines = append(lines, fmt.Sprintf("       Signature: %x", input.ScriptSig))
	//lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	return lines
}

func (output *Output) string() []string {
	var lines []string

	lines = append(lines, fmt.Sprintf("     Output %d:", output.Id))
	lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
	//lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKeyHash))
	lines = append(lines, fmt.Sprintf("       Script: %x", output.ScriptPubKey))

	return lines
}
