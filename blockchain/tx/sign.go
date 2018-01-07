package tx

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/gob"
	"log"
)

func generateSignedInput(input *UnsignedInput, outputs []*Output, privateKey ecdsa.PrivateKey) *SignedInput {
	signatureData := serializeSignatureData(input, outputs)
	r, s, err := ecdsa.Sign(rand.Reader, &privateKey, signatureData)
	if err != nil {
		log.Panic(err)
	}

	signature := append(r.Bytes(), s.Bytes()...)
	return newSignedInput(input, signature)
}

func serializeSignatureData(input *UnsignedInput, outputs []*Output) []byte {
	data := struct {
		Input   *UnsignedInput
		Outputs []*Output
	}{Input: input, Outputs: outputs}

	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	encoder.Encode(data)

	return encoded.Bytes()
}
