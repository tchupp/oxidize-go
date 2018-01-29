package utxo

import (
	"github.com/tclchiam/block_n_go/blockchain/engine/iter"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/identity"
	"github.com/imdario/mergo"
)

type utxoCrawlerEngine struct {
	repository entity.BlockRepository
}

func NewCrawlerEngine(repository entity.BlockRepository) Engine {
	return &utxoCrawlerEngine{repository: repository}
}

func (engine *utxoCrawlerEngine) FindUnspentOutputs(spender *identity.Identity) (*TransactionOutputSet, error) {
	spentOutputs := make(map[string][]*entity.Output)
	outputsForAddress := NewTransactionSet()

	err := iter.ForEachBlock(engine.repository, func(block *entity.Block) {
		for _, transaction := range block.Transactions() {
			mergo.Map(&spentOutputs, findSpentOutputs(transaction, spender))
			outputsForAddress = outputsForAddress.Plus(findOutputsForIdentity(transaction, spender))
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

func findSpentOutputs(transaction *entity.Transaction, spender *identity.Identity) map[string][]*entity.Output {
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
		Filter(func(input *entity.SignedInput) bool { return input.SpentBy(spender) }).
		Reduce(spent, addToUnspent).(map[string][]*entity.Output)
}

func findOutputsForIdentity(transaction *entity.Transaction, identity *identity.Identity) *TransactionOutputSet {
	addToTxSet := func(res interface{}, output *entity.Output) interface{} {
		return res.(*TransactionOutputSet).Add(transaction, output)
	}

	return entity.NewOutputs(transaction.Outputs).
		Filter(func(output *entity.Output) bool { return output.ReceivedBy(identity) }).
		Reduce(NewTransactionSet(), addToTxSet).(*TransactionOutputSet)
}
