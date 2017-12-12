package tx

type Input struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

func (in *Input) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}
