package tx

import (
	"bytes"
	"crypto/sha256"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
)

type TransactionId [sha256.Size]byte

const secretLength = 32

func (txId TransactionId) String() string {
	return hex.EncodeToString(txId[:])
}

func calculateTransactionId(inputs []*UnsignedInput, outputs []*Output) TransactionId {
	return sha256.Sum256(serialize(inputs, outputs))
}

func serialize(inputs []*UnsignedInput, outputs []*Output) []byte {
	data := struct {
		Inputs  []*UnsignedInput
		Outputs []*Output
		Secret  []byte
	}{Inputs: inputs, Outputs: outputs, Secret: generateSecret()}

	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	encoder.Encode(data)

	return encoded.Bytes()
}

func generateSecret() []byte {
	secret := make([]byte, secretLength)
	rand.Read(secret)
	return secret
}
