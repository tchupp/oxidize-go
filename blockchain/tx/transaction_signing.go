package tx

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/tclchiam/block_n_go/crypto"
	"github.com/tclchiam/block_n_go/blockchain/entity"
)

func GenerateSignature(input *entity.UnsignedInput, outputs []*entity.Output, privateKey *crypto.PrivateKey) *crypto.Signature {
	signatureData := serializeSignatureData(input, outputs)
	signature, err := privateKey.Sign(signatureData)
	if err != nil {
		log.Panic(err)
	}

	return signature
}

func serializeSignatureData(input *entity.UnsignedInput, outputs []*entity.Output) []byte {
	data := struct {
		Input   *entity.UnsignedInput
		Outputs []*entity.Output
	}{Input: input, Outputs: outputs}

	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}
