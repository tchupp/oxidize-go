package tx

import (
	"encoding/hex"
)

const subsidy = 10

type (
	TransactionId []byte
	Transaction struct {
		ID        TransactionId
		TxInputs  []*UnsignedInput
		TxOutputs []*Output
	}

	OutputReference struct {
		ID          []byte
		OutputIndex int
	}
)

var EmptyOutputReference = OutputReference{ID: []byte(nil), OutputIndex: -1}

func (txId TransactionId) String() string {
	return hex.EncodeToString(txId)
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.TxInputs) == 1 && !tx.TxInputs[0].isReferencingOutput()
}

func NewGenesisCoinbaseTx(ownerAddress string) *Transaction {
	return NewCoinbaseTx(ownerAddress)
}

func NewCoinbaseTx(minerAddress string) *Transaction {
	input := newCoinbaseTxInput()
	output := NewOutput(subsidy, minerAddress)
	tx := Transaction{nil, []*UnsignedInput{input}, []*Output{output}}
	tx.ID = Hash(&tx)

	return &tx
}

func NewTx(inputs UnsignedInputs, outputs Outputs) *Transaction {
	collectOutputs := func(res interface{}, output *Output) interface{} {
		output.Id = uint(len(res.([]*Output)))
		return append(res.([]*Output), output)
	}

	tx := Transaction{
		TxInputs:  inputs.ToSlice(),
		TxOutputs: outputs.Reduce(make([]*Output, 0), collectOutputs).([]*Output),
	}
	tx.ID = Hash(&tx)

	return &tx
}

func (tx *Transaction) FindOutputsForAddress(address string) *TransactionOutputSet {
	addToTxSet := func(res interface{}, output *Output) interface{} {
		return res.(*TransactionOutputSet).Add(tx, output)
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

	addToTxSet := func(res interface{}, input *UnsignedInput) interface{} {
		transactionId := hex.EncodeToString(input.OutputReference.ID)
		res.(map[string][]uint)[transactionId] = append(res.(map[string][]uint)[transactionId], uint(input.OutputReference.OutputIndex))

		return res
	}

	return tx.Inputs().
		Filter(func(input *UnsignedInput) bool { return input.SpentBy(address) }).
		Reduce(spent, addToTxSet).(map[string][]uint)
}
