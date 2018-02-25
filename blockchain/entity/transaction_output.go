package entity

import (
	"bytes"
	"fmt"

	"github.com/tclchiam/oxidize-go/identity"
)

type Output struct {
	Index         uint32
	Value         uint64
	PublicKeyHash []byte
}

func NewOutput(value uint64, spender *identity.Address) *Output {
	return &Output{Value: value, PublicKeyHash: spender.PublicKeyHash()}
}

func (output *Output) String() string {
	return output.string("")
}

func (output *Output) string(indent string) string {
	return fmt.Sprintf(
		"%sentity.Output{Index: %d, Value: %d, PublicKeyHash: %x}",
		indent,
		output.Index,
		output.Value,
		output.PublicKeyHash,
	)
}

func (output *Output) ReceivedBy(receiver *identity.Address) bool {
	return bytes.Compare(output.PublicKeyHash, receiver.PublicKeyHash()) == 0
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

func EmptyOutputs() Outputs {
	c := make(chan *Output, 0)
	defer close(c)
	return Outputs(c)
}

func NewOutputs(outputs []*Output) Outputs {
	c := make(chan *Output, len(outputs))
	defer close(c)
	for _, output := range outputs {
		c <- output
	}
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
	for output := range outputs {
		res = apply(res, output)
	}
	return res
}

func (outputs Outputs) Add(output *Output) Outputs {
	c := make(chan *Output, len(outputs)+1)
	defer close(c)

	for i := range outputs {
		c <- i
	}
	c <- output
	return Outputs(c)
}

func (outputs Outputs) Append(plus Outputs) Outputs {
	c := make(chan *Output, len(outputs)+len(plus))

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
