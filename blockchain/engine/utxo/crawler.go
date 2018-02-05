package utxo

import (
	"github.com/tclchiam/block_n_go/blockchain/engine/iter"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/identity"
)

type utxoCrawlerEngine struct {
	repository entity.BlockRepository
}

func NewCrawlerEngine(repository entity.BlockRepository) Engine {
	return &utxoCrawlerEngine{repository: repository}
}

func (engine *utxoCrawlerEngine) FindUnspentOutputs(spender *identity.Identity) (*TransactionOutputSet, error) {
	gatherInputs := func(res interface{}, tx *entity.Transaction) interface{} {
		return res.(entity.SignedInputs).Append(entity.NewSignedInputs(tx.Inputs))
	}

	spentOutputs := make(map[*entity.Hash][]*entity.Output)
	outputsForAddress := NewTransactionSet()

	err := iter.ForEachBlock(engine.repository, func(block *entity.Block) {
		spentOutputs = block.Transactions().
			Reduce(entity.EmptySingedInputs(), gatherInputs).(entity.SignedInputs).
			Filter(func(input *entity.SignedInput) bool { return input.SpentBy(spender) }).
			Reduce(spentOutputs, addInputToMap).(map[*entity.Hash][]*entity.Output)

		for _, transaction := range block.Transactions() {
			addToTxSet := func(res interface{}, output *entity.Output) interface{} {
				return res.(*TransactionOutputSet).Add(transaction, output)
			}

			outputsForAddress = entity.NewOutputs(transaction.Outputs).
				Filter(func(output *entity.Output) bool { return output.ReceivedBy(spender) }).
				Reduce(outputsForAddress, addToTxSet).(*TransactionOutputSet)
		}
	})
	if err != nil {
		return nil, err
	}

	return outputsForAddress.
		Filter(isUnspent(spentOutputs)), nil
}

var isUnspent = func(spentOutputs map[*entity.Hash][]*entity.Output) func(*entity.Transaction, *entity.Output) bool {
	return func(transaction *entity.Transaction, output *entity.Output) bool {
		if outputs, ok := spentOutputs[transaction.ID]; ok {
			for _, spentOutput := range outputs {
				if spentOutput.IsEqual(output) {
					return false
				}
			}
		}
		return true
	}
}

var addInputToMap = func(res interface{}, input *entity.SignedInput) interface{} {
	outputs := res.(map[*entity.Hash][]*entity.Output)
	transactionId := input.OutputReference.ID
	outputs[transactionId] = append(outputs[transactionId], input.OutputReference.Output)

	return res
}
