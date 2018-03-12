package txsigning

import (
	log "github.com/sirupsen/logrus"

	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/crypto"
	"github.com/tclchiam/oxidize-go/identity"
	"github.com/tclchiam/oxidize-go/wire"
)

func GenerateSignature(input *entity.UnsignedInput, outputs []*entity.Output, spender *identity.Identity) *crypto.Signature {
	signatureData := serializeSignatureData(input, outputs)
	signature, err := spender.Sign(signatureData)
	if err != nil {
		log.Panic(err)
	}

	return signature
}

func VerifySignature(input *entity.SignedInput, outputs []*entity.Output) (verified bool) {
	unsignedInput := &entity.UnsignedInput{
		PublicKey:       input.PublicKey,
		OutputReference: input.OutputReference,
	}
	signatureData := serializeSignatureData(unsignedInput, outputs)

	return input.PublicKey.Verify(signatureData, input.Signature)
}

func serializeSignatureData(input *entity.UnsignedInput, outputs []*entity.Output) []byte {
	var data []byte

	bytes, err := wire.EncodeUnsignedInput(input)
	if err != nil {
		log.Panic(err)
	}

	data = append(data, bytes...)
	for _, output := range outputs {
		bytes, err := wire.EncodeOutput(output)
		if err != nil {
			log.Panic(err)
		}
		data = append(data, bytes...)
	}

	return data
}
