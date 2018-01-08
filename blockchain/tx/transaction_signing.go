package tx

import (
	"log"

	"github.com/tclchiam/block_n_go/crypto"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/encoding"
)

func GenerateSignature(input *entity.UnsignedInput, outputs []*entity.Output, privateKey *crypto.PrivateKey) *crypto.Signature {
	signatureData := serializeSignatureData(input, outputs, encoding.NewTransactionGobEncoder())
	signature, err := privateKey.Sign(signatureData)
	if err != nil {
		log.Panic(err)
	}

	return signature
}

func serializeSignatureData(input *entity.UnsignedInput, outputs []*entity.Output, encoder entity.TransactionEncoder) []byte {
	var data []byte

	bytes, err := encoder.EncodeUnsignedInput(input)
	if err != nil {
		log.Panic(err)
	}

	data = append(data, bytes...)
	for _, output := range outputs {
		bytes, err := encoder.EncodeOutput(output)
		if err != nil {
			log.Panic(err)
		}
		data = append(data, bytes...)
	}

	return data
}
