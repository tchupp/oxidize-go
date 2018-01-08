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
	spentOutputs := make(map[string][]*entity.Output)
	outputsForAddress := tx.NewTransactionSet()

	err := bc.ForEachBlock(func(block *entity.Block) {
		for _, transaction := range block.Transactions {
			mergo.Map(&spentOutputs, FindSpentOutputs(transaction, address))
			outputsForAddress = outputsForAddress.Plus(FindOutputsForAddress(transaction, address))
		}
	})

	isUnspent := func(transaction *entity.Transaction, output *entity.Output) bool {
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
	sumBalance := func(res interface{}, _ *entity.Transaction, output *entity.Output) interface{} {
		return res.(uint) + output.Value
	}

	return unspentOutputs.Reduce(uint(0), sumBalance).(uint)
}

func FindOutputsForAddress(transaction *entity.Transaction, address string) *tx.TransactionOutputSet {
	addToTxSet := func(res interface{}, output *entity.Output) interface{} {
		return res.(*tx.TransactionOutputSet).Add(transaction, output)
	}

	outputs := tx.NewTransactionSet()
	for _, output := range transaction.Outputs {
		if output.IsLockedWithKey(address) {
			outputs = outputs.Add(transaction, output)
		}
	}

	return entity.NewOutputs(transaction.Outputs).
		Filter(func(output *entity.Output) bool { return output.IsLockedWithKey(address) }).
		Reduce(tx.NewTransactionSet(), addToTxSet).(*tx.TransactionOutputSet)
}

func FindSpentOutputs(transaction *entity.Transaction, address string) map[string][]*entity.Output {
	spent := make(map[string][]*entity.Output)
	if transaction.IsCoinbase() {
		return spent
	}

	addToUnspent := func(res interface{}, input *entity.SignedInput) interface{} {
		transactionId := input.OutputReference.ID.String()
		res.(map[string][]*entity.Output)[transactionId] = append(res.(map[string][]*entity.Output)[transactionId], input.OutputReference.Output)

		return res
	}

	return entity.NewSignedInputs(transaction.Inputs).
		Filter(func(input *entity.SignedInput) bool { return input.SpentBy(address) }).
		Reduce(spent, addToUnspent).(map[string][]*entity.Output)
}
