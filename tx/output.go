package tx

type Output struct {
	Id           int
	Value        int
	ScriptPubKey string
}

func (output *Output) CanBeUnlockedWith(unlockingData string) bool {
	return output.ScriptPubKey == unlockingData
}

type Outputs <-chan Output

func (tx *Transaction) Outputs() Outputs {
	c := make(chan Output, len(tx.TxOutputs))
	go func() {
		for _, output := range tx.TxOutputs {
			c <- output
		}
		close(c)
	}()
	return Outputs(c)
}

func (outputs Outputs) Filter(predicate func(output Output) bool) Outputs {
	outChannel := make(chan Output)

	go func() {
		for output := range outputs {
			if predicate(output) {
				outChannel <- output
			}
		}
		close(outChannel)
	}()
	return Outputs(outChannel)
}

func (outputs Outputs) Plus(plus Outputs) Outputs {
	outChannel := make(chan Output)

	go func() {
		for output := range outputs {
			outChannel <- output
		}
		for output := range plus {
			outChannel <- output
		}
		close(outChannel)
	}()
	return Outputs(outChannel)
}

func (outputs Outputs) ToSlice() []Output {
	slice := make([]Output, 0)
	for o := range outputs {
		slice = append(slice, o)
	}
	return slice
}
