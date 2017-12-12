package tx

type Input struct {
	OutputTransactionId []byte
	OutputId            int
	ScriptSig           string
}

func (input *Input) CanUnlockOutputWith(unlockingData string) bool {
	return input.ScriptSig == unlockingData
}

func (input *Input) isReferencingOutput() bool {
	referencesTransaction := len(input.OutputTransactionId) != 0
	referencesTransactionOutput := input.OutputId != -1

	return referencesTransaction && referencesTransactionOutput
}

func newCoinbaseTxInput(data string) Input {
	return Input{OutputTransactionId: []byte(nil), OutputId: -1, ScriptSig: data}
}
