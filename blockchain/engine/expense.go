package engine

import (
	"fmt"

	"github.com/tclchiam/block_n_go/blockchain/entity/encoding"
	"github.com/tclchiam/block_n_go/blockchain/engine/txsigning"
	"github.com/tclchiam/block_n_go/blockchain/engine/utxo"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/identity"
)

func BuildExpenseTransaction(sender, receiver *identity.Address, expense uint32, engine utxo.Engine) (*entity.Transaction, error) {
	unspentOutputs, err := engine.FindUnspentOutputs(sender)
	if err != nil {
		return nil, err
	}

	balance := calculateBalance(unspentOutputs)
	if balance < expense {
		return nil, fmt.Errorf("account '%s' does not have enough to send '%d', due to balance '%d'", sender, expense, balance)
	}

	finalizedOutputs := entity.EmptyOutputs().
		Add(entity.NewOutput(expense, receiver)).
		Add(entity.NewOutput(balance-expense, sender)).
		Filter(func(output *entity.Output) bool { return output.Value != 0 }).
		Reduce(make([]*entity.Output, 0), collectOutputs).([]*entity.Output)

	signedInputs := unspentOutputs.
		Reduce(entity.EmptyUnsignedInputs(nil), buildInputs(sender)).(entity.UnsignedInputs).
		Reduce(make([]*entity.SignedInput, 0), signInputs(finalizedOutputs, sender)).([]*entity.SignedInput)

	return entity.NewTx(signedInputs, finalizedOutputs, encoding.TransactionProtoEncoder()), nil
}

func buildInputs(sender *identity.Address) func(res interface{}, transaction *entity.Transaction, output *entity.Output) interface{} {
	return func(res interface{}, transaction *entity.Transaction, output *entity.Output) interface{} {
		input := entity.NewUnsignedInput(transaction.ID, output, sender.PublicKey())
		return res.(entity.UnsignedInputs).Add(input)
	}
}

func signInputs(outputs []*entity.Output, sender *identity.Address) func(res interface{}, input *entity.UnsignedInput) interface{} {
	return func(res interface{}, input *entity.UnsignedInput) interface{} {
		signature := txsigning.GenerateSignature(input, outputs, sender, encoding.TransactionProtoEncoder())
		return append(res.([]*entity.SignedInput), entity.NewSignedInput(input, signature))
	}
}

func collectOutputs(res interface{}, output *entity.Output) interface{} {
	outputs := res.([]*entity.Output)
	output.Index = uint32(len(outputs))
	return append(outputs, output)
}
