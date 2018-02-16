package account

import (
	"github.com/tclchiam/oxidize-go/account/utxo"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/identity"
)

func balance(address *identity.Address, engine utxo.Engine) (*Account, error) {
	unspentOutputs, err := engine.FindUnspentOutputs(address)
	if err != nil {
		return nil, err
	}

	unspent := calculateBalance(unspentOutputs)

	return &Account{
		Address:   address,
		Spendable: unspent,
	}, nil
}

func calculateBalance(unspentOutputs *utxo.TransactionOutputSet) uint64 {
	sumBalance := func(res interface{}, _ *entity.Transaction, output *entity.Output) interface{} {
		return res.(uint64) + output.Value
	}

	return unspentOutputs.Reduce(uint64(0), sumBalance).(uint64)
}
