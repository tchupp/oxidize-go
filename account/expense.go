package account

import (
	"fmt"

	"github.com/tclchiam/oxidize-go/blockchain/engine/txsigning"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/blockchain/utxo"
	"github.com/tclchiam/oxidize-go/encoding"
	"github.com/tclchiam/oxidize-go/identity"
)

func buildExpenseTransaction(
	spender *identity.Identity,
	receiver *identity.Address,
	expense uint64,
	spendableOutputs *utxo.OutputSet,
) (*entity.Transaction, error) {
	balance := calculateBalance(spendableOutputs)
	if balance < expense {
		return nil, fmt.Errorf("account '%s' does not have enough to send '%d', due to balance '%d'", spender, expense, balance)
	}

	finalizedOutputs := []*entity.Output{entity.NewOutput(expense, receiver)}
	if balance-expense > 0 {
		finalizedOutputs = append(finalizedOutputs, entity.NewOutput(balance-expense, spender.Address()))
	}

	var signedInputs []*entity.SignedInput
	spendableOutputs.ForEach(func(txId *entity.Hash, output *entity.Output) {
		unsignedInput := entity.NewUnsignedInput(txId, output, spender.PublicKey())
		signature := txsigning.GenerateSignature(unsignedInput, finalizedOutputs, spender, encoding.TransactionProtoEncoder())
		signedInputs = append(signedInputs, entity.NewSignedInput(unsignedInput, signature))
	})

	return entity.NewTx(signedInputs, finalizedOutputs, encoding.TransactionProtoEncoder()), nil
}
