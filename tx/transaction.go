package tx

const subsidy = 10

type (
	Transaction struct {
		ID        TransactionId
		TxInputs  []*UnsignedInput
		TxOutputs []*Output
	}

	OutputReference struct {
		ID          TransactionId
		OutputIndex uint
	}
)

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.TxInputs) == 0
}

func NewGenesisCoinbaseTx(ownerAddress string) *Transaction {
	return NewCoinbaseTx(ownerAddress)
}

func NewCoinbaseTx(minerAddress string) *Transaction {
	var inputs []*UnsignedInput
	outputs := []*Output{NewOutput(subsidy, minerAddress)}

	return newTx(inputs, outputs)
}

func NewTx(inputs UnsignedInputs, outputs Outputs) *Transaction {
	collectOutputs := func(res interface{}, output *Output) interface{} {
		output.Id = uint(len(res.([]*Output)))
		return append(res.([]*Output), output)
	}

	return newTx(inputs.ToSlice(), outputs.Reduce(make([]*Output, 0), collectOutputs).([]*Output))
}

func newTx(inputs []*UnsignedInput, outputs []*Output) *Transaction {
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

	return tx.Outputs().
		Filter(func(output *Output) bool { return output.IsLockedWithKey(address) }).
		Reduce(NewTransactionSet(), addToTxSet).(*TransactionOutputSet)
}

func (tx *Transaction) FindSpentOutputs(address string) map[string][]uint {
	spent := make(map[string][]uint)
	if tx.IsCoinbase() {
		return spent
	}

	addToUnspent := func(res interface{}, input *UnsignedInput) interface{} {
		transactionId := input.OutputReference.ID.String()
		res.(map[string][]uint)[transactionId] = append(res.(map[string][]uint)[transactionId], input.OutputReference.OutputIndex)

		return res
	}

	return tx.Inputs().
		Filter(func(input *UnsignedInput) bool { return input.SpentBy(address) }).
		Reduce(spent, addToUnspent).(map[string][]uint)
}
