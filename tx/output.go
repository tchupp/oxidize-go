package tx

type Output struct {
	Value        int
	ScriptPubKey string
}

func (output *Output) CanBeUnlockedWith(unlockingData string) bool {
	return output.ScriptPubKey == unlockingData
}
