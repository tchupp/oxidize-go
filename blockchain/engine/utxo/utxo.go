package utxo

import (
	"github.com/tclchiam/block_n_go/blockchain/engine/iter"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/blockchain/tx"
	"github.com/tclchiam/block_n_go/storage"
	"github.com/imdario/mergo"
)

type realUTXOEngine struct {
	repository storage.BlockRepository
}

func NewEngine(repository storage.BlockRepository) Engine {
	return &realUTXOEngine{repository: repository}
}

func (engine *realUTXOEngine) FindUnspentOutputs(address string) (*tx.TransactionOutputSet, error) {
	spentOutputs := make(map[string][]*entity.Output)
	outputsForAddress := tx.NewTransactionSet()

	err := iter.ForEachBlock(engine.repository, func(block *entity.Block) {
		for _, transaction := range block.Transactions {
			mergo.Map(&spentOutputs, findSpentOutputs(transaction, address))
			outputsForAddress = outputsForAddress.Plus(findOutputsForAddress(transaction, address))
		}
	})

	return outputsForAddress.Filter(isUnspent(spentOutputs)), err
}

func isUnspent(spentOutputs map[string][]*entity.Output) func(transaction *entity.Transaction, output *entity.Output) bool {
	return func(transaction *entity.Transaction, output *entity.Output) bool {
		if outputs, ok := spentOutputs[transaction.ID.String()]; ok {
			for _, spentOutput := range outputs {
				if spentOutput.IsEqual(output) {
					return false
				}
			}
		}
		return true
	}
}

func findSpentOutputs(transaction *entity.Transaction, address string) map[string][]*entity.Output {
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

func findOutputsForAddress(transaction *entity.Transaction, address string) *tx.TransactionOutputSet {
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
