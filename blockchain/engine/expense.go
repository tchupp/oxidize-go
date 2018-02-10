package engine

import (
	"fmt"

	"github.com/tclchiam/oxidize-go/blockchain/engine/txsigning"
	"github.com/tclchiam/oxidize-go/blockchain/engine/utxo"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/encoding"
	"github.com/tclchiam/oxidize-go/identity"
)

func BuildExpenseTransaction(spender *identity.Identity, receiver *identity.Address, expense uint32, engine utxo.Engine) (*entity.Transaction, error) {
	unspentOutputs, err := engine.FindUnspentOutputs(spender.Address())
	if err != nil {
		return nil, err
	}

	balance := calculateBalance(unspentOutputs)
	if balance < expense {
		return nil, fmt.Errorf("account '%s' does not have enough to send '%d', due to balance '%d'", spender, expense, balance)
	}

	finalizedOutputs := entity.EmptyOutputs().
		Add(entity.NewOutput(expense, receiver)).
		Add(entity.NewOutput(balance-expense, spender.Address())).
		Filter(func(output *entity.Output) bool { return output.Value != 0 }).
		Reduce(make([]*entity.Output, 0), collectOutputs).([]*entity.Output)

	signedInputs := unspentOutputs.
		Reduce(entity.EmptyUnsignedInputs(), buildInputs(spender)).(entity.UnsignedInputs).
		Reduce(make([]*entity.SignedInput, 0), signInputs(finalizedOutputs, spender)).([]*entity.SignedInput)

	return entity.NewTx(signedInputs, finalizedOutputs, encoding.TransactionProtoEncoder()), nil
}

var buildInputs = func(spender *identity.Identity) func(res interface{}, transaction *entity.Transaction, output *entity.Output) interface{} {
	return func(res interface{}, transaction *entity.Transaction, output *entity.Output) interface{} {
		input := entity.NewUnsignedInput(transaction.ID, output, spender.PublicKey())
		return res.(entity.UnsignedInputs).Add(input)
	}
}

var signInputs = func(outputs []*entity.Output, spender *identity.Identity) func(res interface{}, input *entity.UnsignedInput) interface{} {
	return func(res interface{}, input *entity.UnsignedInput) interface{} {
		signature := txsigning.GenerateSignature(input, outputs, spender, encoding.TransactionProtoEncoder())
		return append(res.([]*entity.SignedInput), entity.NewSignedInput(input, signature))
	}
}

var collectOutputs = func(res interface{}, output *entity.Output) interface{} {
	outputs := res.([]*entity.Output)
	output.Index = uint32(len(outputs))
	return append(outputs, output)
}
