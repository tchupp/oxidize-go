package blockchain

import (
	"github.com/tclchiam/block_n_go/tx"
	"github.com/tclchiam/block_n_go/tx/txset"
)

func (bc *Blockchain) findUnspentTransactionOutputs(address string) (unspentOutputs tx.Outputs, err error) {
	spentTransactions := txset.New()

	err = bc.ForEachBlock(func(block *Block) {
		block.ForEachTransaction(func(transaction *tx.Transaction) {
			unspentOutputs = unspentOutputs.Plus(transaction.FindUnspentOutput(spentTransactions, address))
			spentTransactions = transaction.FindSpentOutput(spentTransactions, address)
		})
	})
	return
}

func (bc *Blockchain) ReadBalance(address string) (int, error) {
	unspentOutputs, err := bc.findUnspentTransactionOutputs(address)
	if err != nil {
		return -1, err
	}

	balance := calculateBalance(unspentOutputs)
	return balance, nil
}

func calculateBalance(unspentOutputs tx.Outputs) int {
	balance := 0
	for output := range unspentOutputs {
		balance += output.Value
	}
	return balance
}
