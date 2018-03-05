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

func NewOutput(value uint64, receiver *identity.Address) *Output {
	return &Output{Value: value, PublicKeyHash: receiver.PublicKeyHash()}
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

type Outputs []*Output

func (outputs Outputs) Filter(predicate func(output *Output) bool) Outputs {
	var c []*Output

	for _, output := range outputs {
		if predicate(output) {
			c = append(c, output)
		}
	}
	return Outputs(c)
}

func (outputs Outputs) Add(output *Output) Outputs {
	return append(outputs, output)
}

func (outputs Outputs) Append(plus Outputs) Outputs {
	return append(outputs, plus...)
}

func (outputs Outputs) Contains(output *Output) bool {
	for _, o := range outputs {
		if o.IsEqual(output) {
			return true
		}
	}
	return false
}

func (outputs Outputs) IndexOf(output *Output) int {
	for i, o := range outputs {
		if o.IsEqual(output) {
			return i
		}
	}
	return -1
}

func (outputs Outputs) Remove(output *Output) Outputs {
	index := outputs.IndexOf(output)
	if index == -1 {
		return outputs
	}

	return append(outputs[:index], outputs[index+1:]...)
}
