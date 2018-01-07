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

func NewId(newId []byte) (TransactionId) {
	var id TransactionId
	copy(id[:], newId[:sha256.Size])
	return id
}

func calculateTransactionId(inputs []*SignedInput, outputs []*Output) TransactionId {
	return sha256.Sum256(serializeTxData(inputs, outputs))
}

func serializeTxData(inputs []*SignedInput, outputs []*Output) []byte {
	data := struct {
		Inputs  []*SignedInput
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
