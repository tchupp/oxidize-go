package tx

import (
	"github.com/tclchiam/block_n_go/wallet"
	"bytes"
)

type Output struct {
	Index         uint
	Value         uint
	PublicKeyHash []byte
}

func NewOutput(value uint, address string) *Output {
	publicKeyHash, _ := wallet.AddressToPublicKeyHash(address)

	return &Output{Value: value, PublicKeyHash: publicKeyHash}
}

func (output *Output) IsLockedWithKey(address string) bool {
	publicKeyHash, err := wallet.AddressToPublicKeyHash(address)
	if err != nil {
		return false
	}

	return bytes.Compare(output.PublicKeyHash, publicKeyHash) == 0
}

func (output *Output) IsEqual(other *Output) bool {
	if output == other {
		return true
	}
	if output.Index != other.Index {
		return false
	}
	if output.Value != other.Value {
		return false
	}
	if bytes.Compare(output.PublicKeyHash, other.PublicKeyHash) != 0 {
		return false
	}

	return true
}

type Outputs <-chan *Output

func (tx *Transaction) Outputs() Outputs {
	outputs := tx.TxOutputs
	return NewOutputs(outputs)
}

func EmptyOutputs() Outputs {
	outputs := make([]*Output, 0)
	return NewOutputs(outputs)
}

func NewOutputs(outputs []*Output) Outputs {
	c := make(chan *Output, len(outputs))
	go func() {
		for _, output := range outputs {
			c <- output
		}
		close(c)
	}()
	return Outputs(c)
}

func (outputs Outputs) Filter(predicate func(output *Output) bool) Outputs {
	c := make(chan *Output)

	go func() {
		for output := range outputs {
			if predicate(output) {
				c <- output
			}
		}
		close(c)
	}()
	return Outputs(c)
}

func (outputs Outputs) Reduce(res interface{}, apply func(res interface{}, output *Output) interface{}) interface{} {
	c := make(chan interface{})

	go func() {
		for output := range outputs {
			res = apply(res, output)
		}
		c <- res
	}()
	return <-c
}

func (outputs Outputs) Add(output *Output) Outputs {
	c := make(chan *Output)

	go func() {
		for i := range outputs {
			c <- i
		}
		c <- output
		close(c)
	}()
	return Outputs(c)
}

func (outputs Outputs) Plus(plus Outputs) Outputs {
	c := make(chan *Output)

	go func() {
		for output := range outputs {
			c <- output
		}
		for output := range plus {
			c <- output
		}
		close(c)
	}()
	return Outputs(c)
}

func (outputs Outputs) ToSlice() []*Output {
	slice := make([]*Output, 0)
	for o := range outputs {
		slice = append(slice, o)
	}
	return slice
}
