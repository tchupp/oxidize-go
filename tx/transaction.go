package tx

import (
	"fmt"
)

const subsidy = 10
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type Transaction struct {
	ID      []byte
	Inputs  []Input
	Outputs []Output
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && !tx.Inputs[0].isReferencingOutput()
}

func NewGenesisCoinbaseTx(to string) *Transaction {
	return NewCoinbaseTx(to, genesisCoinbaseData)
}

func NewCoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	input := newCoinbaseTxInput(data)
	output := Output{subsidy, to}
	tx := Transaction{nil, []Input{input}, []Output{output}}
	tx.ID = tx.Hash()

	return &tx
}
