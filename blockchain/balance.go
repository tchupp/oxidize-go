package blockchain

import (
	"github.com/tclchiam/block_n_go/tx"
	"github.com/imdario/mergo"
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
	spentOutputs := make(map[string][]uint)
	outputsForAddress := tx.NewTransactionSet()

	isUnspent := func(transaction *tx.Transaction, output *tx.Output) bool {
		if outputs, ok := spentOutputs[transaction.ID.String()]; ok {
			for _, outputId := range outputs {
				if outputId == output.Id {
					return false
				}
			}
		}
		return true
	}

	err := bc.ForEachBlock(func(block *Block) {
		block.ForEachTransaction(func(transaction *tx.Transaction) {
			mergo.Map(&spentOutputs, transaction.FindSpentOutputs(address))
			outputsForAddress = outputsForAddress.Plus(transaction.FindOutputsForAddress(address))
		})
	})

	unspentOutputs := outputsForAddress.Filter(isUnspent)
	return unspentOutputs, err
}

func calculateBalance(unspentOutputs *tx.TransactionOutputSet) uint {
	sumBalance := func(res interface{}, _ *tx.Transaction, output *tx.Output) interface{} {
		return res.(uint) + output.Value
	}

	return unspentOutputs.Reduce(uint(0), sumBalance).(uint)
}
