package engine

import (
	"github.com/tclchiam/oxidize-go/blockchain/engine/utxo"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/identity"
)

func ReadBalance(identity *identity.Identity, engine utxo.Engine) (uint32, error) {
	unspentOutputs, err := engine.FindUnspentOutputs(identity)
	if err != nil {
		return 0, err
	}

	return calculateBalance(unspentOutputs), nil
}

func calculateBalance(unspentOutputs *utxo.TransactionOutputSet) uint32 {
	sumBalance := func(res interface{}, _ *entity.Transaction, output *entity.Output) interface{} {
		return res.(uint32) + output.Value
	}

	return unspentOutputs.Reduce(uint32(0), sumBalance).(uint32)
}
