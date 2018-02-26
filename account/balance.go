package account

import (
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/blockchain/utxo"
)

func calculateBalance(spendableOutputs *utxo.OutputSet) uint64 {
	balance := uint64(0)
	spendableOutputs.ForEachOutput(func(output *entity.Output) {
		balance += output.Value
	})
	return balance
}
