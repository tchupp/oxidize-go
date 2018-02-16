package utxo

import (
	"github.com/tclchiam/oxidize-go/blockchain/engine/iter"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/identity"
)

type utxoCrawlerEngine struct {
	reader entity.ChainReader
}

func NewCrawlerEngine(reader entity.ChainReader) Engine {
	return &utxoCrawlerEngine{reader: reader}
}

func (engine *utxoCrawlerEngine) FindUnspentOutputs(spender *identity.Address) (*TransactionOutputSet, error) {
	inputs, err := findInputs(engine.reader)
	if err != nil {
		return nil, err
	}

	outputsByTx, err := findOutputsByTransaction(engine.reader)
	if err != nil {
		return nil, err
	}

	spentOutputs := inputs.
		Reduce(make(map[string][]*entity.Output), addInputToMap).(map[string][]*entity.Output)

	return outputsByTx.
		Filter(func(_ *entity.Transaction, output *entity.Output) bool { return output.ReceivedBy(spender) }).
		Filter(isUnspent(spentOutputs)), nil
}

func findInputs(reader entity.ChainReader) (entity.SignedInputs, error) {
	var gatherInputs = func(res interface{}, tx *entity.Transaction) interface{} {
		return res.(entity.SignedInputs).Append(entity.NewSignedInputs(tx.Inputs))
	}

	inputs := entity.EmptySingedInputs()

	err := iter.ForEachBlock(reader, func(block *entity.Block) {
		inputs = block.Transactions().Reduce(inputs, gatherInputs).(entity.SignedInputs)
	})

	return inputs, err
}

func findOutputsByTransaction(reader entity.ChainReader) (*TransactionOutputSet, error) {
	outputsForAddress := NewTransactionSet()

	err := iter.ForEachBlock(reader, func(block *entity.Block) {
		for _, transaction := range block.Transactions() {
			addToTxSet := func(res interface{}, output *entity.Output) interface{} {
				return res.(*TransactionOutputSet).Add(transaction, output)
			}

			outputsForAddress = entity.NewOutputs(transaction.Outputs).
				Reduce(outputsForAddress, addToTxSet).(*TransactionOutputSet)
		}
	})

	return outputsForAddress, err
}

var isUnspent = func(spentOutputs map[string][]*entity.Output) func(*entity.Transaction, *entity.Output) bool {
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

var addInputToMap = func(res interface{}, input *entity.SignedInput) interface{} {
	outputs := res.(map[string][]*entity.Output)
	transactionId := input.OutputReference.ID.String()
	outputs[transactionId] = append(outputs[transactionId], input.OutputReference.Output)

	return res
}
