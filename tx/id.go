package tx

import (
	"crypto/sha256"
	"bytes"
	"encoding/gob"
	"log"
	"encoding/hex"
)

type TransactionId []byte

func (txId TransactionId) String() string {
	return hex.EncodeToString(txId)
}

func calculateTransactionId(inputs []*UnsignedInput, outputs []*Output) TransactionId {
	var hash [32]byte

	hash = sha256.Sum256(serialize(inputs, outputs))

	return hash[:]
}

func serialize(inputs []*UnsignedInput, outputs []*Output) []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)

	data := struct {
		Inputs  []*UnsignedInput
		Outputs []*Output
	}{Inputs: inputs, Outputs: outputs}

	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}
