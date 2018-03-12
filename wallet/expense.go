package wallet

import (
	"fmt"

	"github.com/tclchiam/oxidize-go/blockchain/engine/txsigning"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/identity"
)

func buildExpenseTransaction(
	unspentOutputRefs []*UnspentOutput,
	receiver *identity.Address,
	payback *identity.Address,
	expense uint64,
) (*entity.Transaction, error) {
	balance := uint64(0)
	for _, unspentOutput := range unspentOutputRefs {
		balance += unspentOutput.Output.Value
	}

	if balance < expense {
		return nil, fmt.Errorf("balance %d is not enough to send %d", balance, expense)
	}

	finalizedOutputs := []*entity.Output{entity.NewOutput(expense, receiver)}
	if balance-expense > 0 {
		finalizedOutputs = append(finalizedOutputs, entity.NewOutput(balance-expense, payback))
	}

	var signedInputs []*entity.SignedInput
	for _, outputRef := range unspentOutputRefs {
		unsignedInput := entity.NewUnsignedInput(outputRef.TxId, outputRef.Output, outputRef.Identity.PublicKey())
		signature := txsigning.GenerateSignature(unsignedInput, finalizedOutputs, outputRef.Identity)
		signedInputs = append(signedInputs, entity.NewSignedInput(unsignedInput, signature))
	}
	return entity.NewTx(signedInputs, finalizedOutputs), nil
}
