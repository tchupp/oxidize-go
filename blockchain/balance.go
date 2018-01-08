package blockchain

import (
	"github.com/imdario/mergo"
	"github.com/tclchiam/block_n_go/blockchain/tx"
	"github.com/tclchiam/block_n_go/blockchain/entity"
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

	err := bc.ForEachBlock(func(block *entity.Block) {
		block.ForEachTransaction(func(transaction *tx.Transaction) {
			mergo.Map(&spentOutputs, FindSpentOutputs(transaction, address))
			outputsForAddress = outputsForAddress.Plus(FindOutputsForAddress(transaction, address))
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

func FindOutputsForAddress(transaction *tx.Transaction, address string) *tx.TransactionOutputSet {
	addToTxSet := func(res interface{}, output *tx.Output) interface{} {
		return res.(*tx.TransactionOutputSet).Add(transaction, output)
	}

	outputs := tx.NewTransactionSet()
	for _, output := range transaction.TxOutputs {
		if output.IsLockedWithKey(address) {
			outputs = outputs.Add(transaction, output)
		}
	}

	return tx.NewOutputs(transaction.TxOutputs).
		Filter(func(output *tx.Output) bool { return output.IsLockedWithKey(address) }).
		Reduce(tx.NewTransactionSet(), addToTxSet).(*tx.TransactionOutputSet)
}

func FindSpentOutputs(transaction *tx.Transaction, address string) map[string][]*tx.Output {
	spent := make(map[string][]*tx.Output)
	if transaction.IsCoinbase() {
		return spent
	}

	addToUnspent := func(res interface{}, input *tx.SignedInput) interface{} {
		transactionId := input.OutputReference.ID.String()
		res.(map[string][]*tx.Output)[transactionId] = append(res.(map[string][]*tx.Output)[transactionId], input.OutputReference.Output)

		return res
	}

	return transaction.Inputs().
		Filter(func(input *tx.SignedInput) bool { return input.SpentBy(address) }).
		Reduce(spent, addToUnspent).(map[string][]*tx.Output)
}
