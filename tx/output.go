package tx

type Output struct {
	Id           int
	Value        int
	ScriptPubKey string
}

func NewOutput(value int, to string) *Output {
	return &Output{Value: value, ScriptPubKey: to}
}

func (output *Output) CanBeUnlockedWith(unlockingData string) bool {
	return output.ScriptPubKey == unlockingData
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
