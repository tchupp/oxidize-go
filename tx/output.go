package tx

type Output struct {
	Value        int
	ScriptPubKey string
}

func (out *Output) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}
