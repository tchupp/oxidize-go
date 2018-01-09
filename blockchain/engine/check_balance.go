package engine

import (
	"github.com/tclchiam/block_n_go/blockchain/engine/utxo"
	"github.com/tclchiam/block_n_go/blockchain/entity"
)

func ReadBalance(address string, engine utxo.Engine) (uint, error) {
	unspentOutputs, err := engine.FindUnspentOutputs(address)
	if err != nil {
		return 0, err
	}

	return calculateBalance(unspentOutputs), nil
}

func calculateBalance(unspentOutputs *utxo.TransactionOutputSet) uint {
	sumBalance := func(res interface{}, _ *entity.Transaction, output *entity.Output) interface{} {
		return res.(uint) + output.Value
	}

	return unspentOutputs.Reduce(uint(0), sumBalance).(uint)
}
