package tx

import (
	"encoding/hex"
	"fmt"

	"github.com/tclchiam/block_n_go/tx/txset"
)

const subsidy = 10
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type Transaction struct {
	ID        []byte
	TxInputs  []Input
	TxOutputs []Output
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
	output := Output{0, subsidy, to}
	tx := Transaction{nil, []Input{input}, []Output{output}}
	tx.ID = tx.Hash()

	return &tx
}

func (tx *Transaction) FindUnspentOutput(spent *txset.TransactionSet, address string) []Output {
	transactionId := hex.EncodeToString(tx.ID)

	return tx.Outputs().
		Filter(func(output Output) bool { return !spent.Contains(transactionId, output.Id) }).
		Filter(func(output Output) bool { return output.CanBeUnlockedWith(address) }).
		ToSlice()
}

func (tx *Transaction) FindSpentOutput(spent *txset.TransactionSet, address string) *txset.TransactionSet {
	if tx.IsCoinbase() {
		return spent
	}

	addToTxSet := func(res interface{}, input Input) interface{} {
		transactionId := hex.EncodeToString(input.OutputTransactionId)
		return res.(*txset.TransactionSet).Add(transactionId, input.OutputId)
	}

	return tx.Inputs().
		Filter(func(input Input) bool { return input.CanUnlockOutputWith(address) }).
		Reduce(spent, addToTxSet).(*txset.TransactionSet)
}
