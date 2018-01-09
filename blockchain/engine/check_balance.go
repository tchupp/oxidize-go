package engine

import (
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/blockchain/engine/utxo"
	"github.com/tclchiam/block_n_go/blockchain/tx"
	"github.com/tclchiam/block_n_go/storage"
)

func ReadBalance(address string, repository storage.BlockRepository) (uint, error) {
	unspentOutputs, err := utxo.NewEngine(repository).FindUnspentOutputs(address)
	if err != nil {
		return 0, err
	}

	return calculateBalance(unspentOutputs), nil
}

func calculateBalance(unspentOutputs *tx.TransactionOutputSet) uint {
	sumBalance := func(res interface{}, _ *entity.Transaction, output *entity.Output) interface{} {
		return res.(uint) + output.Value
	}

	return unspentOutputs.Reduce(uint(0), sumBalance).(uint)
}
