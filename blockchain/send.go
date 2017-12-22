package blockchain

import (
	"fmt"
	"encoding/hex"
	"log"

	"github.com/tclchiam/block_n_go/tx"
	"github.com/tclchiam/block_n_go/wallet"
)

func (bc *Blockchain) buildExpenseTransaction(sender, receiver *wallet.Wallet, expense uint) (*tx.Transaction, error) {
	senderAddress := sender.GetAddress()

	unspentOutputs, err := bc.findUnspentOutputs(senderAddress)
	if err != nil {
		return nil, err
	}

	balance := calculateBalance(unspentOutputs)
	if balance < expense {
		return nil, fmt.Errorf("account '%x' does not have enough to send '%d', due to balance '%d'", senderAddress, expense, balance)
	}

	liquidBalance := uint(0)
	takeMinimumToMeetExpense := func(_ string, output *tx.Output) bool {
		take := liquidBalance < expense
		if take {
			liquidBalance += output.Value
		}
		return take
	}

	buildInputs := func(res interface{}, transactionId string, output *tx.Output) interface{} {
		b, err := hex.DecodeString(transactionId)
		if err != nil {
			log.Panic(err)
		}

		input := tx.NewInput(b, int(output.Id), sender.PublicKey)
		return res.(tx.Inputs).Add(input)
	}

	inputs := unspentOutputs.
		Filter(takeMinimumToMeetExpense).
		Reduce(tx.NewInputs(nil), buildInputs).(tx.Inputs)

	receiverAddress := receiver.GetAddress()
	outputs := tx.EmptyOutputs().
		Add(tx.NewOutput(expense, receiverAddress))

	if liquidBalance-expense > 0 {
		outputs = outputs.Add(tx.NewOutput(liquidBalance-expense, senderAddress))
	}

	return tx.NewTx(inputs, outputs), nil
}
