package blockchain

import (
	"github.com/imdario/mergo"
	"github.com/tclchiam/block_n_go/blockchain/tx"
	"github.com/tclchiam/block_n_go/blockchain/block"
)

func (bc *Blockchain) ReadBalance(address string) (uint, error) {
	unspentOutputs, err := bc.findUnspentOutputs(address)
	if err != nil {
		return 0, err
	}

	balance := calculateBalance(unspentOutputs)
	return balance, nil
}

func (bc *Blockchain) findUnspentOutputs(address string) (*tx.TransactionOutputSet, error) {
	spentOutputs := make(map[string][]*tx.Output)
	outputsForAddress := tx.NewTransactionSet()

	err := bc.ForEachBlock(func(block *block.Block) {
		block.ForEachTransaction(func(transaction *tx.Transaction) {
			mergo.Map(&spentOutputs, transaction.FindSpentOutputs(address))
			outputsForAddress = outputsForAddress.Plus(transaction.FindOutputsForAddress(address))
		})
	})

	isUnspent := func(transaction *tx.Transaction, output *tx.Output) bool {
		if outputs, ok := spentOutputs[transaction.ID.String()]; ok {
			for _, spentOutput := range outputs {
				if spentOutput.IsEqual(output) {
					return false
				}
			}
		}
		return true
	}

	return outputsForAddress.Filter(isUnspent), err
}

func calculateBalance(unspentOutputs *tx.TransactionOutputSet) uint {
	sumBalance := func(res interface{}, _ *tx.Transaction, output *tx.Output) interface{} {
		return res.(uint) + output.Value
	}

	return unspentOutputs.Reduce(uint(0), sumBalance).(uint)
}
