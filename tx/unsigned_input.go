package tx

import (
	"bytes"

	"github.com/tclchiam/block_n_go/wallet"
	"crypto/rand"
)

type UnsignedInput struct {
	OutputReference OutputReference
	PublicKey       []byte
}

func NewUnsignedInput(outputTransactionId TransactionId, outputId uint, senderPublicKey []byte) *UnsignedInput {
	reference := OutputReference{ID: outputTransactionId, OutputIndex: outputId}

	return newUnsignedInput(reference, senderPublicKey)
}

func newUnsignedInput(reference OutputReference, senderPublicKey []byte) *UnsignedInput {
	return &UnsignedInput{
		OutputReference: reference,
		PublicKey:       senderPublicKey,
	}
}

func (input *UnsignedInput) SpentBy(address string) bool {
	publicKeyHash, err := wallet.AddressToPublicKeyHash(address)
	if err != nil {
		return false
	}

	lockingHash := wallet.HashPubKey(input.PublicKey)

	return bytes.Compare(lockingHash, publicKeyHash) == 0
}

func (input *UnsignedInput) isReferencingOutput() bool {
	referencesTransaction := len(input.OutputReference.ID) != 0
	referencesTransactionOutput := input.OutputReference.OutputIndex != 123456789

	return referencesTransaction && referencesTransactionOutput
}

func newCoinbaseTxInput() *UnsignedInput {
	randData := make([]byte, 20)
	rand.Read(randData)
	return newUnsignedInput(EmptyOutputReference, randData)
}

type UnsignedInputs <-chan *UnsignedInput

func (tx *Transaction) Inputs() UnsignedInputs {
	c := make(chan *UnsignedInput, len(tx.TxInputs))
	go func() {
		for _, input := range tx.TxInputs {
			c <- input
		}
		close(c)
	}()
	return UnsignedInputs(c)
}

func NewInputs(inputs []*UnsignedInput) UnsignedInputs {
	c := make(chan *UnsignedInput, len(inputs))
	go func() {
		for _, input := range inputs {
			c <- input
		}
		close(c)
	}()
	return UnsignedInputs(c)
}

func (inputs UnsignedInputs) Filter(predicate func(input *UnsignedInput) bool) UnsignedInputs {
	c := make(chan *UnsignedInput)

	go func() {
		for input := range inputs {
			if predicate(input) {
				c <- input
			}
		}
		close(c)
	}()
	return UnsignedInputs(c)
}

func (inputs UnsignedInputs) Reduce(res interface{}, apply func(res interface{}, input *UnsignedInput) interface{}) interface{} {
	c := make(chan interface{})

	go func() {
		for input := range inputs {
			res = apply(res, input)
		}
		c <- res
	}()
	return <-c
}

func (inputs UnsignedInputs) Add(input *UnsignedInput) UnsignedInputs {
	c := make(chan *UnsignedInput)

	go func() {
		for i := range inputs {
			c <- i
		}
		c <- input
		close(c)
	}()
	return UnsignedInputs(c)
}

func (inputs UnsignedInputs) ToSlice() []*UnsignedInput {
	slice := make([]*UnsignedInput, 0)
	for i := range inputs {
		slice = append(slice, i)
	}
	return slice
}
