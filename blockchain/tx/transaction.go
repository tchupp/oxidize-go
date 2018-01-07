package tx

import (
	"crypto/ecdsa"
)

const subsidy = 10

type (
	Transaction struct {
		ID        TransactionId
		TxInputs  []*SignedInput
		TxOutputs []*Output
	}

	OutputReference struct {
		ID     TransactionId
		Output *Output
	}
)

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.TxInputs) == 0
}

func NewGenesisCoinbaseTx(ownerAddress string) *Transaction {
	return NewCoinbaseTx(ownerAddress)
}

func NewCoinbaseTx(minerAddress string) *Transaction {
	var inputs []*SignedInput
	outputs := []*Output{NewOutput(subsidy, minerAddress)}

	return newTx(inputs, outputs)
}

func collectOutputs(res interface{}, output *Output) interface{} {
	outputs := res.([]*Output)
	output.Index = uint(len(outputs))
	return append(outputs, output)
}

func NewTx(inputs *UnsignedInputs, outputs *Outputs, privateKey ecdsa.PrivateKey) *Transaction {
	finalizedOutputs := outputs.Reduce(make([]*Output, 0), collectOutputs).([]*Output)

	signInputs := func(res interface{}, input *UnsignedInput) interface{} {
		signedInput := generateSignedInput(input, finalizedOutputs, privateKey)
		return append(res.([]*SignedInput), signedInput)
	}
	signedInputs := inputs.Reduce(make([]*SignedInput, 0), signInputs).([]*SignedInput)

	return newTx(signedInputs, finalizedOutputs)
}

func newTx(inputs []*SignedInput, outputs []*Output) *Transaction {
	return &Transaction{
		ID:        calculateTransactionId(inputs, outputs),
		TxInputs:  inputs,
		TxOutputs: outputs,
	}
}

func (tx *Transaction) FindOutputsForAddress(address string) *TransactionOutputSet {
	addToTxSet := func(res interface{}, output *Output) interface{} {
		return res.(*TransactionOutputSet).Add(tx, output)
	}

	outputs := NewTransactionSet()
	for _, output := range tx.TxOutputs {
		if output.IsLockedWithKey(address) {
			outputs = outputs.Add(tx, output)
		}
	}

	return tx.Outputs().
		Filter(func(output *Output) bool { return output.IsLockedWithKey(address) }).
		Reduce(NewTransactionSet(), addToTxSet).(*TransactionOutputSet)
}

func (tx *Transaction) FindSpentOutputs(address string) map[string][]*Output {
	spent := make(map[string][]*Output)
	if tx.IsCoinbase() {
		return spent
	}

	addToUnspent := func(res interface{}, input *SignedInput) interface{} {
		transactionId := input.OutputReference.ID.String()
		res.(map[string][]*Output)[transactionId] = append(res.(map[string][]*Output)[transactionId], input.OutputReference.Output)

		return res
	}

	return tx.Inputs().
		Filter(func(input *SignedInput) bool { return input.SpentBy(address) }).
		Reduce(spent, addToUnspent).(map[string][]*Output)
}
