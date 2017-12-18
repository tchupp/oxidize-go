package tx

import (
	"fmt"
	"encoding/hex"
)

const subsidy = 10
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type Transaction struct {
	ID        []byte
	TxInputs  []*Input
	TxOutputs []*Output
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.TxInputs) == 1 && !tx.TxInputs[0].isReferencingOutput()
}

func NewGenesisCoinbaseTx(to string) *Transaction {
	return NewCoinbaseTx(to, genesisCoinbaseData)
}

func NewCoinbaseTx(to string, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	input := newCoinbaseTxInput(data)
	output := NewOutput(subsidy, to)
	tx := Transaction{nil, []*Input{input}, []*Output{output}}
	tx.ID = tx.Hash()

	return &tx
}

func NewTx(inputs Inputs, outputs Outputs) *Transaction {
	collectOutputs := func(res interface{}, output *Output) interface{} {
		output.Id = len(res.([]*Output))
		return append(res.([]*Output), output)
	}

	tx := Transaction{
		TxInputs:  inputs.ToSlice(),
		TxOutputs: outputs.Reduce(make([]*Output, 0), collectOutputs).([]*Output),
	}
	tx.ID = tx.Hash()

	return &tx
}

func (tx *Transaction) FindOutputsForAddress(address string) *TransactionSet {
	transactionId := hex.EncodeToString(tx.ID)

	addToTxSet := func(res interface{}, output *Output) interface{} {
		return res.(*TransactionSet).Add(transactionId, output)
	}

	return tx.Outputs().
		Filter(func(output *Output) bool { return output.CanBeUnlockedWith(address)}).
		Reduce(NewTransactionSet(), addToTxSet).(*TransactionSet)
}

func (tx *Transaction) FindSpentOutputs(address string) map[string][]int {
	spent := make(map[string][]int)
	if tx.IsCoinbase() {
		return spent
	}

	addToTxSet := func(res interface{}, input *Input) interface{} {
		transactionId := hex.EncodeToString(input.OutputTransactionId)
		res.(map[string][]int)[transactionId] = append(res.(map[string][]int)[transactionId], input.OutputId)

		return res
	}

	return tx.Inputs().
		Filter(func(input *Input) bool { return input.CanUnlockOutputWith(address) }).
		Reduce(spent, addToTxSet).(map[string][]int)
}
