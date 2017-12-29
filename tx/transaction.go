package tx

import (
	"encoding/hex"
)

const subsidy = 10

type (
	Transaction struct {
		ID        []byte
		TxInputs  []*Input
		TxOutputs []*Output
	}
)

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.TxInputs) == 1 && !tx.TxInputs[0].isReferencingOutput()
}

func NewGenesisCoinbaseTx(ownerAddress string) *Transaction {
	return NewCoinbaseTx(ownerAddress)
}

func NewCoinbaseTx(minerAddress string) *Transaction {
	input := newCoinbaseTxInput()
	output := NewOutput(subsidy, minerAddress)
	tx := Transaction{nil, []*Input{input}, []*Output{output}}
	tx.ID = tx.Hash()

	return &tx
}

func NewTx(inputs Inputs, outputs Outputs) *Transaction {
	collectOutputs := func(res interface{}, output *Output) interface{} {
		output.Id = uint(len(res.([]*Output)))
		return append(res.([]*Output), output)
	}

	tx := Transaction{
		TxInputs:  inputs.ToSlice(),
		TxOutputs: outputs.Reduce(make([]*Output, 0), collectOutputs).([]*Output),
	}
	tx.ID = tx.Hash()

	return &tx
}

func (tx *Transaction) FindOutputsForAddress(address string) *TransactionOutputSet {
	transactionId := hex.EncodeToString(tx.ID)

	addToTxSet := func(res interface{}, output *Output) interface{} {
		return res.(*TransactionOutputSet).Add(transactionId, output)
	}

	return tx.Outputs().
		Filter(func(output *Output) bool { return output.IsLockedWithKey(address) }).
		Reduce(NewTransactionSet(), addToTxSet).(*TransactionOutputSet)
}

func (tx *Transaction) FindSpentOutputs(address string) map[string][]uint {
	spent := make(map[string][]uint)
	if tx.IsCoinbase() {
		return spent
	}

	addToTxSet := func(res interface{}, input *Input) interface{} {
		transactionId := hex.EncodeToString(input.OutputTransactionId)
		res.(map[string][]uint)[transactionId] = append(res.(map[string][]uint)[transactionId], uint(input.OutputId))

		return res
	}

	return tx.Inputs().
		Filter(func(input *Input) bool { return input.SpentBy(address) }).
		Reduce(spent, addToTxSet).(map[string][]uint)
}
