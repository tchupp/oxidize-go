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
