package blockchain

import (
	"github.com/tclchiam/block_n_go/tx"
	"github.com/tclchiam/block_n_go/tx/txset"
)

func (bc *Blockchain) FindUnspentTransactionOutputs(address string) (unspentOutputs tx.Outputs, err error) {
	spentTransactions := txset.New()

	err = bc.ForEachBlock(func(block *Block) {
		block.ForEachTransaction(func(transaction *tx.Transaction) {
			unspentOutputs = unspentOutputs.Plus(transaction.FindUnspentOutput(spentTransactions, address))
			spentTransactions = transaction.FindSpentOutput(spentTransactions, address)
		})
	})
	return
}
