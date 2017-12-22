package tx

import (
	"bytes"
	"github.com/tclchiam/block_n_go/wallet"
	"crypto/rand"
)

type Input struct {
	OutputTransactionId []byte
	OutputId            int
	Signature           []byte
	PublicKey           []byte
}

func NewInput(outputTransactionId []byte, outputId int, senderPublicKey []byte) *Input {
	return &Input{
		OutputTransactionId: outputTransactionId,
		OutputId:            outputId,
		Signature:           nil,
		PublicKey:           senderPublicKey,
	}
}

func (input *Input) SpentBy(address string) bool {
	publicKeyHash, err := wallet.AddressToPublicKeyHash(address)
	if err != nil {
		return false
	}

	lockingHash := wallet.HashPubKey(input.PublicKey)

	return bytes.Compare(lockingHash, publicKeyHash) == 0
}

func (input *Input) isReferencingOutput() bool {
	referencesTransaction := len(input.OutputTransactionId) != 0
	referencesTransactionOutput := input.OutputId != -1

	return referencesTransaction && referencesTransactionOutput
}

func newCoinbaseTxInput() *Input {
	randData := make([]byte, 20)
	rand.Read(randData)
	return NewInput([]byte(nil), -1, randData)
}

type Inputs <-chan *Input

func (tx *Transaction) Inputs() Inputs {
	c := make(chan *Input, len(tx.TxInputs))
	go func() {
		for _, input := range tx.TxInputs {
			c <- input
		}
		close(c)
	}()
	return Inputs(c)
}

func NewInputs(inputs []*Input) Inputs {
	c := make(chan *Input, len(inputs))
	go func() {
		for _, input := range inputs {
			c <- input
		}
		close(c)
	}()
	return Inputs(c)
}

func (inputs Inputs) Filter(predicate func(input *Input) bool) Inputs {
	c := make(chan *Input)

	go func() {
		for input := range inputs {
			if predicate(input) {
				c <- input
			}
		}
		close(c)
	}()
	return Inputs(c)
}

func (inputs Inputs) Reduce(res interface{}, apply func(res interface{}, input *Input) interface{}) interface{} {
	c := make(chan interface{})

	go func() {
		for input := range inputs {
			res = apply(res, input)
		}
		c <- res
	}()
	return <-c
}

func (inputs Inputs) Add(input *Input) Inputs {
	c := make(chan *Input)

	go func() {
		for i := range inputs {
			c <- i
		}
		c <- input
		close(c)
	}()
	return Inputs(c)
}

func (inputs Inputs) ToSlice() []*Input {
	slice := make([]*Input, 0)
	for i := range inputs {
		slice = append(slice, i)
	}
	return slice
}
